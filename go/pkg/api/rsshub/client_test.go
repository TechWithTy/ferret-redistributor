package rsshub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoutes(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	routes, err := client.GetRoutes(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "ok", routes.Status)
	assert.Contains(t, routes.Data, "bilibili")
}

func TestFetchFeed(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.FetchFeed(context.Background(), FeedRequest{
		Path: "/bilibili/fav/123",
		Query: map[string]string{
			"limit": "5",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "application/rss+xml", resp.ContentType)
	assert.Contains(t, string(resp.Body), "<rss")
}

func TestForceRefresh(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	result, err := client.ForceRefresh(context.Background(), ForceRefreshRequest{
		TargetURL: "https://rsshub.app/test",
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", result.Status)
}

func TestSearchRadar(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	result, err := client.SearchRadar(context.Background(), RadarSearchRequest{
		URL: "https://bilibili.com",
	})
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	assert.Equal(t, "/bilibili/fav/:uid", result.Results[0].Path)
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	version, err := client.GetVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", version.Version)
}

func TestFetchFeedMissingPath(t *testing.T) {
	t.Parallel()

	client := NewClient()
	_, err := client.FetchFeed(context.Background(), FeedRequest{})
	assert.ErrorIs(t, err, ErrMissingRoutePath)
}

func TestForceRefreshMissingURL(t *testing.T) {
	t.Parallel()

	client := NewClient()
	_, err := client.ForceRefresh(context.Background(), ForceRefreshRequest{})
	assert.ErrorIs(t, err, ErrMissingURL)
}

func TestFetchDebugJSON(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	body, err := client.FetchDebugJSON(context.Background(), "/debug/test", map[string]string{"lang": "en"})
	require.NoError(t, err)
	assert.Contains(t, string(body), `"mode":"debug.json"`)
}

func TestFetchItemDebugHTML(t *testing.T) {
	t.Parallel()

	server := newFakeServer(t)
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	body, err := client.FetchItemDebugHTML(context.Background(), "/debug/test", 2, nil)
	require.NoError(t, err)
	assert.Contains(t, string(body), "2.debug.html")
}

func newFakeServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/routes", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, RoutesResponse{
			Status: "ok",
			Data: map[string]RoutesGroupData{
				"bilibili": {
					Name: "Bilibili",
					Routes: []RouteDetail{
						{Path: "/bilibili/fav/:uid"},
					},
				},
			},
		})
	})

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, VersionResponse{Version: "1.0.0", Commit: "abc123"})
	})

	mux.HandleFunc("/api/radar/search", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, RadarSearchResponse{
			Results: []RadarSearchItem{
				{Title: "Favorites", Path: "/bilibili/fav/:uid"},
			},
		})
	})

	mux.HandleFunc("/api/force-refresh/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, ForceRefreshResponse{Status: "ok", Message: "refreshing"})
	})

	mux.HandleFunc("/bilibili/fav/123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "5", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte("<rss><channel></channel></rss>"))
	})

	mux.HandleFunc("/debug/test", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("format") {
		case "debug.json":
			require.Equal(t, "en", r.URL.Query().Get("lang"))
			writeJSON(t, w, map[string]string{"mode": "debug.json"})
		case "2.debug.html":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<div>2.debug.html</div>"))
		default:
			http.NotFound(w, r)
		}
	})

	server := httptest.NewServer(http.TimeoutHandler(mux, time.Second, "timeout"))
	return server
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(v))
}
