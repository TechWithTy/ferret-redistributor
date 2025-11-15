package instagram

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "strings"
)

// WebhookHandler handles Instagram Graph webhooks for comments.
// It supports GET verification and POST delivery with signature validation.
type WebhookHandler struct {
    AppSecret   string
    VerifyToken string
    // OnComment is called for each comment change delivered in the payload.
    OnComment func(ctx context.Context, c CommentChange) error
    // OnMessage is called for each message event (requires IG Messaging permissions).
    OnMessage func(ctx context.Context, m MessageChange) error
    // OnMention is called when your account is mentioned.
    OnMention func(ctx context.Context, m MentionChange) error
    // OnInteraction receives a generic interaction event (comment, mention, like).
    OnInteraction func(ctx context.Context, i InteractionChange) error
}

// CommentChange represents a single Instagram comment event.
type CommentChange struct {
    CommentID string
    MediaID   string
    Text      string
    FromID    string
    Username  string
}

// MessageChange represents a single Instagram messaging event.
type MessageChange struct {
    MessageID string
    FromID    string
    Username  string
    Text      string
    Timestamp int64
    ThreadID  string
}

// MentionChange represents a mention event of your account.
type MentionChange struct {
    MentionID string
    MediaID   string
    Text      string
    FromID    string
    Username  string
}

// InteractionChange is a generic wrapper for interactions such as comments, mentions, and likes.
type InteractionChange struct {
    Type      string // comment | mention | like
    ID        string
    MediaID   string
    Text      string
    FromID    string
    Username  string
}

// ServeHTTP implements http.Handler.
func (h WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.handleVerify(w, r)
    case http.MethodPost:
        h.handleDelivery(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (h WebhookHandler) handleVerify(w http.ResponseWriter, r *http.Request) {
    mode := r.URL.Query().Get("hub.mode")
    token := r.URL.Query().Get("hub.verify_token")
    challenge := r.URL.Query().Get("hub.challenge")
    if mode == "subscribe" && token == h.VerifyToken {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(challenge))
        return
    }
    w.WriteHeader(http.StatusForbidden)
}

func (h WebhookHandler) handleDelivery(w http.ResponseWriter, r *http.Request) {
    // Validate signature if app secret is set
    var body []byte
    var err error
    if h.AppSecret != "" {
        sig := r.Header.Get("X-Hub-Signature-256")
        body, err = io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        if !validateSignature(h.AppSecret, sig, body) {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
    } else {
        // If no app secret provided, read body once
        body, err = io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    }

    var payload graphWebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Only process instagram object (IG Graph)
    if strings.ToLower(payload.Object) != "instagram" {
        w.WriteHeader(http.StatusOK)
        return
    }

    ctx := r.Context()
    if h.OnComment != nil {
        for _, entry := range payload.Entry {
            for _, ch := range entry.Changes {
                switch strings.ToLower(ch.Field) {
                case "comments":
                    var v igCommentValue
                    if err := json.Unmarshal(ch.Value, &v); err != nil { continue }
                    cc := CommentChange{CommentID: v.ID, MediaID: v.MediaID, Text: v.Text, FromID: v.From.ID, Username: v.From.Username}
                    if h.OnComment != nil { _ = h.OnComment(ctx, cc) }
                    if h.OnInteraction != nil { _ = h.OnInteraction(ctx, InteractionChange{Type: "comment", ID: cc.CommentID, MediaID: cc.MediaID, Text: cc.Text, FromID: cc.FromID, Username: cc.Username}) }
                case "mentions":
                    var v igMentionValue
                    if err := json.Unmarshal(ch.Value, &v); err != nil { continue }
                    mc := MentionChange{MentionID: v.ID, MediaID: v.MediaID, Text: v.Text, FromID: v.From.ID, Username: v.From.Username}
                    if h.OnMention != nil { _ = h.OnMention(ctx, mc) }
                    if h.OnInteraction != nil { _ = h.OnInteraction(ctx, InteractionChange{Type: "mention", ID: mc.MentionID, MediaID: mc.MediaID, Text: mc.Text, FromID: mc.FromID, Username: mc.Username}) }
                case "messages":
                    var v igMessageValue
                    if err := json.Unmarshal(ch.Value, &v); err != nil { continue }
                    me := MessageChange{MessageID: v.ID, FromID: v.From.ID, Username: v.From.Username, Text: v.Text, Timestamp: v.Timestamp, ThreadID: v.ThreadID}
                    if h.OnMessage != nil { _ = h.OnMessage(ctx, me) }
                    if h.OnInteraction != nil { _ = h.OnInteraction(ctx, InteractionChange{Type: "message", ID: me.MessageID, Text: me.Text, FromID: me.FromID, Username: me.Username}) }
                case "likes":
                    var v igLikeValue
                    if err := json.Unmarshal(ch.Value, &v); err != nil { continue }
                    if h.OnInteraction != nil { _ = h.OnInteraction(ctx, InteractionChange{Type: "like", ID: v.ID, MediaID: v.MediaID, FromID: v.From.ID, Username: v.From.Username}) }
                }
            }
        }
    }
    w.WriteHeader(http.StatusOK)
}

func validateSignature(appSecret, headerSig string, body []byte) bool {
    // Header format: "sha256=HEX"
    parts := strings.SplitN(headerSig, "=", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "sha256" {
        return false
    }
    mac := hmac.New(sha256.New, []byte(appSecret))
    mac.Write(body)
    expected := mac.Sum(nil)
    given, err := hex.DecodeString(parts[1])
    if err != nil {
        return false
    }
    return hmac.Equal(expected, given)
}

// Graph webhook payloads for Instagram object
type graphWebhookPayload struct {
    Object string         `json:"object"`
    Entry  []graphEntry   `json:"entry"`
}

type graphEntry struct {
    ID      string           `json:"id"`
    Time    int64            `json:"time"`
    Changes []graphChange    `json:"changes"`
}

type graphChange struct {
    Field string          `json:"field"`
    Value json.RawMessage `json:"value"`
}

// Expected shape for comments value
type igCommentValue struct {
    ID      string `json:"id"`
    MediaID string `json:"media_id"`
    Text    string `json:"text"`
    From    struct {
        ID       string `json:"id"`
        Username string `json:"username"`
    } `json:"from"`
}

// Expected shape for mentions value
type igMentionValue struct {
    ID      string `json:"id"`
    MediaID string `json:"media_id"`
    Text    string `json:"text"`
    From    struct {
        ID       string `json:"id"`
        Username string `json:"username"`
    } `json:"from"`
}

// Expected shape for messages value
type igMessageValue struct {
    ID        string `json:"id"`
    Text      string `json:"text"`
    Timestamp int64  `json:"timestamp"`
    ThreadID  string `json:"thread_id"`
    From      struct {
        ID       string `json:"id"`
        Username string `json:"username"`
    } `json:"from"`
}

// Expected shape for likes value (if delivered)
type igLikeValue struct {
    ID      string `json:"id"`
    MediaID string `json:"media_id"`
    From    struct {
        ID       string `json:"id"`
        Username string `json:"username"`
    } `json:"from"`
}
