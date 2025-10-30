package recurpost

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	rp "github.com/bitesinbyte/ferret/pkg/api/recurpost"
	"github.com/joho/godotenv"
)

// helper to read required creds
func getCreds(t *testing.T) (email, pass string) {
	t.Helper()
	// Load root .env so tests can use RECURPOST_EMAIL/RECURPOST_PASSWORD from project root
	// From this package path (pkg/api/recurpost/tests) to go root is four levels up
	_ = godotenv.Load("../../../../.env")
	email = os.Getenv("RECURPOST_EMAIL")
	pass = os.Getenv("RECURPOST_PASSWORD")
	if email == "" || pass == "" {
		t.Skip("RECURPOST_EMAIL/RECURPOST_PASSWORD not set; skipping RecurPost client tests")
	}
	return
}

func newMockServer(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for p, h := range routes {
		mux.HandleFunc(p, h)
	}
	return httptest.NewServer(mux)
}

func newClient(base string) *rp.Client {
	return rp.NewClient(rp.WithBaseURL(base))
}

func TestUserLogin(t *testing.T) {
	email, pass := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/user_login": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			var body struct {
				EmailID string `json:"emailid"`
				PassKey string `json:"pass_key"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			if body.EmailID != email || body.PassKey != pass {
				http.Error(w, "creds", http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "access_token": "token123"})
		},
	})
	defer srv.Close()

	cli := newClient(srv.URL)
	out, err := cli.UserLogin.Login(context.Background(), rp.UserLoginRequest{EmailID: email, PassKey: pass})
	if err != nil {
		t.Fatalf("login error: %v", err)
	}
	if !out.Success || out.AccessToken == "" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestConnectSocialAccountURLs(t *testing.T) {
	email, pass := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/connect_social_account_urls": func(w http.ResponseWriter, r *http.Request) {
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["emailid"] != email || body["pass_key"] != pass {
				http.Error(w, "creds", http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"urls": map[string]string{"facebook": "https://fb.example/connect"}})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL)
	out, err := cli.SocialConnect.GetURLs(context.Background(), rp.ConnectSocialAccountURLsRequest{EmailID: email, PassKey: pass})
	if err != nil {
		t.Fatalf("connect urls error: %v", err)
	}
	if out.URLs["facebook"] == "" {
		t.Fatalf("expected facebook url, got: %+v", out)
	}
}

func TestSocialAccountList(t *testing.T) {
	email, _ := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/social_account_list": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("emailid"); got != email {
				http.Error(w, "missing email", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"accounts": []map[string]string{{"id": "acc1", "provider": "twitter", "handle": "@x", "display_name": "X"}}})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL)
	out, err := cli.SocialAccounts.List(context.Background(), rp.SocialAccountListRequest{EmailID: email})
	if err != nil {
		t.Fatalf("accounts list error: %v", err)
	}
	if len(out.Accounts) == 0 {
		t.Fatal("expected at least one account")
	}
}

func TestLibraryListAndAddContent(t *testing.T) {
	email, pass := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/library_list": func(w http.ResponseWriter, r *http.Request) {
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["emailid"] != email || body["pass_key"] != pass {
				http.Error(w, "creds", http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"libraries": []map[string]string{{"id": "lib1", "name": "Default"}}})
		},
		"/api/add_content_in_library": func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["emailid"] != email || body["pass_key"] != pass {
				http.Error(w, "creds", http.StatusUnauthorized)
				return
			}
			if strings.TrimSpace(body["id"].(string)) == "" || strings.TrimSpace(body["message"].(string)) == "" {
				http.Error(w, "bad", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "content_id": "content123"})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL)
	libs, err := cli.Libraries.List(context.Background(), rp.LibraryListRequest{EmailID: email, PassKey: pass})
	if err != nil {
		t.Fatalf("library list error: %v", err)
	}
	if len(libs.Libraries) == 0 {
		t.Fatal("expected libraries")
	}
	add, err := cli.Libraries.AddContent(context.Background(), rp.AddContentInLibraryRequest{EmailID: email, PassKey: pass, ID: libs.Libraries[0].ID, Message: "hello world"})
	if err != nil {
		t.Fatalf("add content error: %v", err)
	}
	if !add.Success || add.ContentID == "" {
		t.Fatalf("unexpected add response: %+v", add)
	}
}

func TestHistoryData(t *testing.T) {
	email, pass := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/history_data": func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []string{"ok"}})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL)
	out, err := cli.History.Data(context.Background(), rp.HistoryDataRequest{EmailID: email, PassKey: pass, ID: "acc1", IsGetVideoUpdates: "true"})
	if err != nil {
		t.Fatalf("history error: %v", err)
	}
	if len(out.Data) == 0 {
		t.Fatal("expected history data")
	}
}

func TestPublishAndAI(t *testing.T) {
	email, pass := getCreds(t)
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/post_content": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "post_id": "p1"})
		},
		"/api/generate_content_with_ai": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"text": "generated"}})
		},
		"/api/generate_image_with_ai": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"image_url": "https://example.com/img.png"}})
		},
	})
	defer srv.Close()
	cli := newClient(srv.URL)
	pub, err := cli.Publishing.Post(context.Background(), rp.PostContentRequest{EmailID: email, PassKey: pass, ID: "acc1", Message: "hi"})
	if err != nil {
		t.Fatalf("post content error: %v", err)
	}
	if !pub.Success || pub.PostID == "" {
		t.Fatalf("unexpected post response: %+v", pub)
	}
	gc, err := cli.AI.GenerateContent(context.Background(), rp.GenerateContentWithAIRequest{EmailID: email, PassKey: pass, PromptText: "topic"})
	if err != nil {
		t.Fatalf("gen content error: %v", err)
	}
	if len(gc.Data) == 0 {
		t.Fatal("expected ai content data")
	}
	gi, err := cli.AI.GenerateImage(context.Background(), rp.GenerateImageWithAIRequest{EmailID: email, PassKey: pass, PromptText: "image"})
	if err != nil {
		t.Fatalf("gen image error: %v", err)
	}
	if len(gi.Data) == 0 {
		t.Fatal("expected ai image data")
	}
}

func TestErrorHandling(t *testing.T) {
	email, pass := getCreds(t)
	// Mock endpoints returning 400, 401, 405
	srv := newMockServer(t, map[string]http.HandlerFunc{
		"/api/user_login": func(w http.ResponseWriter, r *http.Request) {
			// Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": "unauthorized", "message": "invalid creds"})
		},
		"/api/add_content_in_library": func(w http.ResponseWriter, r *http.Request) {
			// Validation error
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": "validation_error", "message": "missing message"})
		},
		"/api/history_data": func(w http.ResponseWriter, r *http.Request) {
			// Method not allowed
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": "method_not_allowed", "message": "use POST"})
		},
	})
	defer srv.Close()

	cli := newClient(srv.URL)
	// 401
	_, err := cli.UserLogin.Login(context.Background(), rp.UserLoginRequest{EmailID: email, PassKey: pass})
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	if !rp.IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized, got: %v", err)
	}

	// 400
	_, err = cli.Libraries.AddContent(context.Background(), rp.AddContentInLibraryRequest{EmailID: email, PassKey: pass, ID: "lib1", Message: ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var apiErr *rp.APIError
	if ok := errors.As(err, &apiErr); !ok || apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 APIError, got: %T %v", err, err)
	}

	// 405
	_, err = cli.History.Data(context.Background(), rp.HistoryDataRequest{EmailID: email, PassKey: pass})
	if err == nil {
		t.Fatal("expected method not allowed error")
	}
	if ok := errors.As(err, &apiErr); !ok || apiErr.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 APIError, got: %T %v", err, err)
	}
}
