package instagram_test

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
)

func TestWebhookVerifyAndComment(t *testing.T) {
    called := make(chan struct{}, 1)
    h := ig.WebhookHandler{
        AppSecret:   "secret",
        VerifyToken: "verify",
        OnComment: func(_ context.Context, c ig.CommentChange) error {
            if c.Text == "hi" && c.MediaID == "m" && c.FromID == "u" { called <- struct{}{} }
            return nil
        },
    }

    // Verify
    req := httptest.NewRequest(http.MethodGet, "/?hub.mode=subscribe&hub.verify_token=verify&hub.challenge=abc", nil)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK || w.Body.String() != "abc" { t.Fatalf("verify failed: %d %q", w.Code, w.Body.String()) }

    // Delivery
    body := map[string]any{
        "object": "instagram",
        "entry": []any{
            map[string]any{
                "id": "1", "time": 1,
                "changes": []any{
                    map[string]any{
                        "field": "comments",
                        "value": map[string]any{
                            "id": "c", "media_id": "m", "text": "hi",
                            "from": map[string]any{"id": "u", "username": "x"},
                        },
                    },
                },
            },
        },
    }
    b, _ := json.Marshal(body)
    mac := hmac.New(sha256.New, []byte("secret"))
    mac.Write(b)
    sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
    preq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
    preq.Header.Set("X-Hub-Signature-256", sig)
    pw := httptest.NewRecorder()
    h.ServeHTTP(pw, preq)
    if pw.Code != http.StatusOK { t.Fatalf("post code: %d", pw.Code) }
    select {
    case <-called:
    default:
        t.Fatal("OnComment not called")
    }
}

