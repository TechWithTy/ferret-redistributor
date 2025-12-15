package groupme_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/stretchr/testify/require"
)

func TestParseAndValidateWebhook_OK_QueryToken(t *testing.T) {
	body := mustReadFile(t, "testdata/webhook_message.json")
	r := httptest.NewRequest(http.MethodPost, "/webhooks/groupme?token=secret", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	ev, err := groupme.ParseAndValidateWebhook(r, groupme.WebhookConfig{Token: "secret"})
	require.NoError(t, err)
	require.Equal(t, "m_123", ev.MessageID)
	require.Equal(t, "g_123", ev.GroupID)
	require.Equal(t, "u_123", ev.SenderID)
	require.Equal(t, "user", ev.SenderType)
	require.Equal(t, "!ping", ev.Text)
}

func TestParseAndValidateWebhook_OK_HeaderToken(t *testing.T) {
	body := mustReadFile(t, "testdata/webhook_message.json")
	r := httptest.NewRequest(http.MethodPost, "/webhooks/groupme", bytes.NewReader(body))
	r.Header.Set("X-Webhook-Token", "secret")

	ev, err := groupme.ParseAndValidateWebhook(r, groupme.WebhookConfig{Token: "secret"})
	require.NoError(t, err)
	require.Equal(t, "m_123", ev.MessageID)
}

func TestParseAndValidateWebhook_Unauthorized(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/webhooks/groupme?token=wrong", bytes.NewReader([]byte(`{}`)))
	_, err := groupme.ParseAndValidateWebhook(r, groupme.WebhookConfig{Token: "secret"})
	require.Error(t, err)
	require.ErrorIs(t, err, groupme.ErrUnauthorized)
}

func TestParseAndValidateWebhook_BadJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/webhooks/groupme?token=secret", bytes.NewReader([]byte("{")))
	_, err := groupme.ParseAndValidateWebhook(r, groupme.WebhookConfig{Token: "secret"})
	require.Error(t, err)
	require.ErrorIs(t, err, groupme.ErrBadRequest)
}

func TestParseAndValidateWebhook_BadShape(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/webhooks/groupme?token=secret", bytes.NewReader([]byte(`{"id":"x"}`)))
	_, err := groupme.ParseAndValidateWebhook(r, groupme.WebhookConfig{Token: "secret"})
	require.Error(t, err)
	require.ErrorIs(t, err, groupme.ErrBadRequest)
}

func mustReadFile(t *testing.T, rel string) []byte {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	// When tests run from this package, cwd is go/pkg/external/groupme
	p := filepath.Join(wd, rel)
	b, err := os.ReadFile(p)
	require.NoError(t, err)
	return b
}
