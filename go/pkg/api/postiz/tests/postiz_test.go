package postiz_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	pz "github.com/bitesinbyte/ferret/pkg/api/postiz"
	"github.com/joho/godotenv"
)

func getAPIKey(t *testing.T) string {
	t.Helper()
	_ = godotenv.Load("../../../../.env")
	k := os.Getenv("POSTIZ_API_KEY")
	if k == "" {
		t.Skip("POSTIZ_API_KEY not set")
	}
	return k
}

func newMock(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for p, h := range routes {
		mux.HandleFunc(p, h)
	}
	return httptest.NewServer(mux)
}

func newClient(base, key string) *pz.Client {
	return pz.NewClient(pz.WithBaseURL(base), pz.WithAPIKey(key))
}

func TestIntegrationsAndFindSlot(t *testing.T) {
	key := getAPIKey(t)
	srv := newMock(t, map[string]http.HandlerFunc{
		"/integrations": func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != key {
				http.Error(w, "auth", http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "i1", "name": "Name", "identifier": "facebook"}})
		},
		"/find-slot/i1": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]string{"date": "2025-01-01T10:00:00.000Z"})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL, key)
	ints, err := cli.Integrations.List(context.Background())
	if err != nil || len(ints) == 0 {
		t.Fatalf("integrations: %v %v", ints, err)
	}
	slot, err := cli.Slots.Find(context.Background(), "i1")
	if err != nil || slot.Date == "" {
		t.Fatalf("slot: %+v %v", slot, err)
	}
}

func TestUpload(t *testing.T) {
	key := getAPIKey(t)
	srv := newMock(t, map[string]http.HandlerFunc{
		"/upload": func(w http.ResponseWriter, r *http.Request) {
			mr, err := r.MultipartReader()
			if err != nil {
				http.Error(w, "multipart", http.StatusBadRequest)
				return
			}
			part, _ := mr.NextPart()
			if part == nil || part.FormName() != "file" {
				http.Error(w, "file", http.StatusBadRequest)
				return
			}
			b, _ := io.ReadAll(part)
			if len(b) == 0 {
				http.Error(w, "empty", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "u1", "name": "file.png", "path": "https://uploads/file.png"})
		},
		"/upload-from-url": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "u2", "name": "file.png", "path": "https://uploads/file.png"})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL, key)
	up, err := cli.Upload.Upload(context.Background(), "file.png", []byte("data"))
	if err != nil || up.ID == "" {
		t.Fatalf("upload: %+v %v", up, err)
	}
	up2, err := cli.Upload.UploadFromURL(context.Background(), pz.UploadFromURLRequest{URL: "https://example.com/a.png"})
	if err != nil || up2.ID == "" {
		t.Fatalf("upload-from-url: %+v %v", up2, err)
	}
}

func TestPostsAndVideo(t *testing.T) {
	key := getAPIKey(t)
	srv := newMock(t, map[string]http.HandlerFunc{
		"/posts": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(map[string]any{"posts": []map[string]any{{"id": "p1", "content": "hi", "publishDate": "2024-12-09T05:06:00.000Z", "state": "DRAFT", "integration": map[string]string{"id": "i1"}}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode([]map[string]string{{"postId": "p2", "integration": "i1"}})
			}
		},
		"/posts/p2": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "p2"})
		},
		"/generate-video": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode([]map[string]string{{"id": "v1", "path": "https://uploads/vid.mp4"}})
		},
		"/video/function": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"voices": []map[string]string{{"id": "v1", "name": "voice"}}})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL, key)

	list, err := cli.Posts.List(context.Background(), pz.PostsListRequest{StartDate: "2024-12-01T00:00:00.000Z", EndDate: "2024-12-31T23:59:59.999Z"})
	if err != nil || len(list.Posts) == 0 {
		t.Fatalf("list: %+v %v", list, err)
	}
	res, err := cli.Posts.CreateOrUpdate(context.Background(), pz.CreateUpdatePostRequest{Type: "draft", Posts: []pz.PostInput{{Integration: pz.IntegrationInput{ID: "i1"}, Value: []pz.PostContent{{Content: "hello"}}}}})
	if err != nil || len(res) == 0 {
		t.Fatalf("create: %+v %v", res, err)
	}
	del, err := cli.Posts.Delete(context.Background(), "p2")
	if err != nil || del.ID == "" {
		t.Fatalf("delete: %+v %v", del, err)
	}

	vids, err := cli.Video.Generate(context.Background(), pz.GenerateVideoRequest{Type: "image-text-slides", Output: "vertical", CustomParams: map[string]any{"voice": "abc"}})
	if err != nil || len(vids) == 0 {
		t.Fatalf("gen video: %+v %v", vids, err)
	}
	fn, err := cli.Video.Function(context.Background(), pz.VideoFunctionRequest{FunctionName: "loadVoices", Identifier: "image-text-slides"})
	if err != nil || len(fn.Voices) == 0 {
		t.Fatalf("fn: %+v %v", fn, err)
	}
}
