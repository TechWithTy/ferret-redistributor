package groupme_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/stretchr/testify/require"
)

func TestUserClient_ListBots_AddsTokenQueryParam(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v3/bots", r.URL.Path)
		require.Equal(t, "t123", r.URL.Query().Get("token"))
		require.NotEmpty(t, r.URL.Query().Get("page"))
		require.NotEmpty(t, r.URL.Query().Get("per_page"))
		if r.URL.Query().Get("page") == "1" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"response": []map[string]any{
					{"bot_id": "b1", "group_id": "g1", "name": "Bot", "callback_url": "https://x", "avatar_url": "https://a"},
				},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"response": []map[string]any{}})
	}))
	t.Cleanup(srv.Close)

	c, err := groupme.NewUserClient(groupme.UserConfig{BaseURL: srv.URL + "/v3", AccessToken: "t123"})
	require.NoError(t, err)

	bots, err := c.ListBots(context.Background())
	require.NoError(t, err)
	require.Len(t, bots, 1)
	require.Equal(t, "b1", bots[0].BotID)
	require.Equal(t, "g1", bots[0].GroupID)
}

func TestUserClient_ListBots_Paginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]any{"response": []map[string]any{{"bot_id": "b1", "group_id": "g1"}}})
		case "2":
			_ = json.NewEncoder(w).Encode(map[string]any{"response": []map[string]any{{"bot_id": "b2", "group_id": "g1"}}})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"response": []map[string]any{}})
		}
	}))
	t.Cleanup(srv.Close)

	c, err := groupme.NewUserClient(groupme.UserConfig{BaseURL: srv.URL + "/v3", AccessToken: "t123"})
	require.NoError(t, err)

	bots, err := c.ListBots(context.Background())
	require.NoError(t, err)
	require.Len(t, bots, 2)
	require.Equal(t, "b1", bots[0].BotID)
	require.Equal(t, "b2", bots[1].BotID)
}

func TestUserClient_ListGroups_NormalizesMembersCount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v3/groups", r.URL.Path)
		require.Equal(t, "t123", r.URL.Query().Get("token"))
		if r.URL.Query().Get("page") == "1" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"response": []map[string]any{
					{"id": "g1", "name": "G", "creator_user_id": "u1", "members": []map[string]any{{"user_id": "u1"}, {"user_id": "u2"}}},
				},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"response": []map[string]any{}})
	}))
	t.Cleanup(srv.Close)

	c, err := groupme.NewUserClient(groupme.UserConfig{BaseURL: srv.URL + "/v3", AccessToken: "t123"})
	require.NoError(t, err)

	groups, err := c.ListGroups(context.Background())
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, 2, groups[0].MembersCount)
}
