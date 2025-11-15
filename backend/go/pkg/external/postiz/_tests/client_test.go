package postiz_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	postiz "github.com/bitesinbyte/ferret/pkg/external/postiz"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *postiz.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := postiz.New(postiz.Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	return client
}

func requireAuthHeader(t *testing.T, r *http.Request) {
	t.Helper()
	require.Equal(t, "test-key", r.Header.Get("Authorization"))
}

func TestClient_ListIntegrations(t *testing.T) {
	sample := []postiz.Integration{
		{ID: "int_1", Name: "Demo", Identifier: "facebook"},
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/integrations", r.URL.Path)
		requireAuthHeader(t, r)
		_ = json.NewEncoder(w).Encode(sample)
	})

	res, err := client.ListIntegrations(context.Background())
	require.NoError(t, err)
	require.Equal(t, sample, res)
}

func TestClient_FindNextSlot(t *testing.T) {
	date := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/find-slot/int_123", r.URL.Path)
		requireAuthHeader(t, r)
		_ = json.NewEncoder(w).Encode(postiz.SlotResponse{Date: date})
	})

	res, err := client.FindNextSlot(context.Background(), "int_123")
	require.NoError(t, err)
	require.Equal(t, date, res)
}

func TestClient_ListPosts(t *testing.T) {
	start := time.Unix(0, 0)
	end := start.Add(24 * time.Hour)
	expected := postiz.PostsResponse{Posts: []postiz.Post{{ID: "post_1"}}}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/posts", r.URL.Path)
		require.Equal(t, start.UTC().Format(time.RFC3339), r.URL.Query().Get("startDate"))
		require.Equal(t, end.UTC().Format(time.RFC3339), r.URL.Query().Get("endDate"))
		require.Equal(t, "cust_1", r.URL.Query().Get("customer"))
		requireAuthHeader(t, r)
		_ = json.NewEncoder(w).Encode(expected)
	})

	res, err := client.ListPosts(context.Background(), postiz.ListPostsParams{
		StartDate: start,
		EndDate:   end,
		Customer:  "cust_1",
	})
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestClient_CreateOrUpdatePosts(t *testing.T) {
	payload := postiz.CreatePostsRequest{
		Type:      "schedule",
		ShortLink: true,
		Date:      time.Now(),
		Posts: []postiz.PostDraft{
			{
				Integration: postiz.IntegrationRef{ID: "int_1"},
				Value:       []postiz.PostValue{{Content: "hello"}},
			},
		},
	}
	expected := []postiz.CreatePostsResponse{{PostID: "post_1", Integration: "int_1"}}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/posts", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		requireAuthHeader(t, r)

		body, _ := io.ReadAll(r.Body)
		require.Contains(t, string(body), `"schedule"`)
		_ = json.NewEncoder(w).Encode(expected)
	})

	res, err := client.CreateOrUpdatePosts(context.Background(), payload)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestClient_DeletePost(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/posts/post_123", r.URL.Path)
		requireAuthHeader(t, r)
		_ = json.NewEncoder(w).Encode(postiz.DeletePostResponse{ID: "post_123"})
	})

	res, err := client.DeletePost(context.Background(), "post_123")
	require.NoError(t, err)
	require.Equal(t, "post_123", res.ID)
}

func TestClient_UploadFromURL(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/upload-from-url", r.URL.Path)
		requireAuthHeader(t, r)
		var req postiz.UploadFromURLRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.Equal(t, "https://example.com/image.png", req.URL)
		_ = json.NewEncoder(w).Encode(postiz.FileAsset{ID: "file_1"})
	})

	res, err := client.UploadFromURL(context.Background(), "https://example.com/image.png")
	require.NoError(t, err)
	require.Equal(t, "file_1", res.ID)
}

func TestClient_UploadFile(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/upload", r.URL.Path)
		requireAuthHeader(t, r)
		require.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		require.NoError(t, r.ParseMultipartForm(10<<20))
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		data, _ := io.ReadAll(file)
		require.Equal(t, "demo.txt", header.Filename)
		require.Equal(t, "content", string(data))
		_ = json.NewEncoder(w).Encode(postiz.FileAsset{ID: "file_upload"})
	})

	res, err := client.UploadFile(context.Background(), "demo.txt", strings.NewReader("content"))
	require.NoError(t, err)
	require.Equal(t, "file_upload", res.ID)
}

func TestClient_GenerateVideo(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/generate-video", r.URL.Path)
		requireAuthHeader(t, r)
		_ = json.NewEncoder(w).Encode([]postiz.VideoAsset{{ID: "vid_1"}})
	})

	payload := postiz.VideoRequest{Type: "image-text-slides", Output: "vertical"}
	res, err := client.GenerateVideo(context.Background(), payload)
	require.NoError(t, err)
	require.Equal(t, "vid_1", res[0].ID)
}

func TestClient_LoadVideoVoices(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/video/function", r.URL.Path)
		requireAuthHeader(t, r)
		var req postiz.VideoFunctionRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.Equal(t, "loadVoices", req.FunctionName)
		_ = json.NewEncoder(w).Encode(postiz.VideoFunctionResponse{
			Voices: []postiz.VideoVoice{{ID: "voice_1", Name: "Demo"}},
		})
	})

	payload := postiz.VideoFunctionRequest{FunctionName: "loadVoices", Identifier: "image-text-slides"}
	res, err := client.LoadVideoVoices(context.Background(), payload)
	require.NoError(t, err)
	require.Len(t, res.Voices, 1)
}
