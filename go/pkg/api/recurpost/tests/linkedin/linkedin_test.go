package recurpost_test

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    rp "github.com/bitesinbyte/ferret/pkg/api/recurpost"
    "github.com/joho/godotenv"
)

// TestLinkedinPost verifies posting to LinkedIn via RecurPost's /api/post_content flow.
// It uses local helpers so this file can run alone if needed.
func TestLinkedinPost_Mock(t *testing.T) {
    email, pass := getCredsLinkedIn(t)

    // Mock RecurPost backend for /api/post_content
    srv := newMockServerLinkedIn(t, map[string]http.HandlerFunc{
        "/api/post_content": func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodPost {
                http.Error(w, "method", http.StatusMethodNotAllowed)
                return
            }
            var body rp.PostContentRequest
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, "bad json", http.StatusBadRequest)
                return
            }
            // Basic assertions: required fields + LinkedIn-specific message present
            if body.EmailID != email || body.PassKey != pass || body.ID == "" || body.Message == "" {
                http.Error(w, "missing required fields", http.StatusBadRequest)
                return
            }
            if body.LNMessage == "" {
                http.Error(w, "missing ln_message", http.StatusBadRequest)
                return
            }
            _ = json.NewEncoder(w).Encode(map[string]any{"success": true, "post_id": "ln_123"})
        },
    })
    defer srv.Close()

    cli := newClientLI(srv.URL)
    in := rp.PostContentRequest{
        EmailID:   email,
        PassKey:   pass,
        ID:        "linkedin_acc_1",
        Message:   "fallback message",
        LNMessage: "linkedin-specific message",
    }
    out, err := cli.Publishing.Post(context.Background(), in)
    if err != nil {
        t.Fatalf("post linkedin error: %v", err)
    }
    if !out.Success || out.PostID == "" {
        t.Fatalf("unexpected response: %+v", out)
    }
}

// Local helpers so this file can run standalone if invoked directly.
func getCredsLinkedIn(t *testing.T) (email, pass string) {
    t.Helper()
    _ = godotenv.Load("../../../../.env")
    email = os.Getenv("RECURPOST_EMAIL")
    pass = os.Getenv("RECURPOST_PASSWORD")
    if email == "" || pass == "" {
        t.Skip("RECURPOST_EMAIL/RECURPOST_PASSWORD not set; skipping")
    }
    return
}

func newMockServerLinkedIn(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
    t.Helper()
    mux := http.NewServeMux()
    for p, h := range routes {
        mux.HandleFunc(p, h)
    }
    return httptest.NewServer(mux)
}

// newClientLI constructs a RecurPost client pointing to a base URL (e.g., mock server).
func newClientLI(base string) *rp.Client {
    return rp.NewClient(rp.WithBaseURL(base))
}
