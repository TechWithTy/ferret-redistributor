package groupme

// MessageCallback is the payload GroupMe sends to a bot callback URL.
// We model only the fields we need for local validation + simple replies.
//
// Example fields are documented by GroupMe's bot callback payload.
type MessageCallback struct {
	ID         string `json:"id"`
	GroupID    string `json:"group_id"`
	SenderID   string `json:"sender_id"`
	SenderType string `json:"sender_type"` // "user" | "bot"
	Name       string `json:"name"`
	Text       string `json:"text"`
	System     bool   `json:"system"`
	CreatedAt  int64  `json:"created_at"`
}

// WebhookEvent is a normalized representation of a callback.
type WebhookEvent struct {
	MessageID  string
	GroupID    string
	SenderID   string
	SenderType string
	Name       string
	Text       string
	System     bool
}
