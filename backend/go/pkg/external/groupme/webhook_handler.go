package groupme

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	// ErrUnauthorized is returned when the shared secret is missing or invalid.
	ErrUnauthorized = errors.New("groupme: unauthorized")
	// ErrBadRequest is returned when the payload is invalid.
	ErrBadRequest = errors.New("groupme: bad request")
)

// WebhookConfig configures webhook validation/behavior.
type WebhookConfig struct {
	// Token is a shared secret used to validate inbound callbacks.
	// Provide it either as query string `?token=...` or header `X-Webhook-Token: ...`.
	Token string
}

// ParseAndValidateWebhook reads and validates a GroupMe callback request.
//
// Validation:
// - Shared secret token must match (query `token` or header `X-Webhook-Token`).
// - JSON must parse and include a minimally valid shape.
//
// Behavior:
// - Bot-sent messages are not rejected; callers can decide to ignore them to prevent loops.
func ParseAndValidateWebhook(r *http.Request, cfg WebhookConfig) (WebhookEvent, error) {
	if strings.TrimSpace(cfg.Token) == "" {
		return WebhookEvent{}, fmt.Errorf("%w: GROUPME_WEBHOOK_TOKEN not set", ErrUnauthorized)
	}
	if !validToken(r, cfg.Token) {
		return WebhookEvent{}, ErrUnauthorized
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return WebhookEvent{}, ErrBadRequest
	}
	defer r.Body.Close()

	var payload MessageCallback
	if err := json.Unmarshal(body, &payload); err != nil {
		return WebhookEvent{}, ErrBadRequest
	}

	// Minimal shape validation
	if strings.TrimSpace(payload.ID) == "" ||
		strings.TrimSpace(payload.GroupID) == "" ||
		strings.TrimSpace(payload.SenderType) == "" ||
		strings.TrimSpace(payload.SenderID) == "" {
		return WebhookEvent{}, ErrBadRequest
	}

	return WebhookEvent{
		MessageID:  payload.ID,
		GroupID:    payload.GroupID,
		SenderID:   payload.SenderID,
		SenderType: strings.ToLower(strings.TrimSpace(payload.SenderType)),
		Name:       payload.Name,
		Text:       payload.Text,
		System:     payload.System,
	}, nil
}

func validToken(r *http.Request, expected string) bool {
	exp := strings.TrimSpace(expected)
	if exp == "" {
		return false
	}
	if got := strings.TrimSpace(r.URL.Query().Get("token")); got != "" {
		return got == exp
	}
	if got := strings.TrimSpace(r.Header.Get("X-Webhook-Token")); got != "" {
		return got == exp
	}
	if got := strings.TrimSpace(r.Header.Get("Authorization")); got != "" {
		// Allow: Authorization: Bearer <token>
		lower := strings.ToLower(got)
		if strings.HasPrefix(lower, "bearer ") {
			return strings.TrimSpace(got[len("Bearer "):]) == exp
		}
	}
	return false
}
