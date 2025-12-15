package groupme_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/stretchr/testify/require"
)

func TestClient_PostBotMessage_OK(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v3/bots/post", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)

	c := groupme.New(groupme.Config{BaseURL: srv.URL + "/v3", HTTPTimeout: 2 * time.Second})
	err := c.PostBotMessage(context.Background(), "bot_1", "hello")
	require.NoError(t, err)
	require.Equal(t, "bot_1", got["bot_id"])
	require.Equal(t, "hello", got["text"])
}

func TestClient_PostBotMessage_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	t.Cleanup(srv.Close)

	c := groupme.New(groupme.Config{BaseURL: srv.URL + "/v3"})
	err := c.PostBotMessage(context.Background(), "bot_1", "hello")
	require.Error(t, err)
}
