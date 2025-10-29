package recurpost_test

import (
    "context"
    "os"
    "testing"
    "time"
    "strings"

    rp "github.com/bitesinbyte/ferret/pkg/api/recurpost"
)

// Integration test against the live RecurPost API for LinkedIn posting.
// This test is DISABLED by default. Enable by setting RUN_RECURPOST_INTEGRATION=1
// and provide credentials + account id via environment:
//   RECURPOST_BASE_URL (e.g., https://api.recurpost.com)
//   RECURPOST_EMAIL
//   RECURPOST_PASSWORD
//   RECURPOST_LINKEDIN_ACCOUNT_ID
// Optional:
//   RECURPOST_LINKEDIN_TEST_MESSAGE (default auto-generated)
//   RECURPOST_SCHEDULE_MINS (if set, schedules N minutes in the future)
func TestLinkedinPost_Live(t *testing.T) {
    if !envBool("RUN_RECURPOST_INTEGRATION") {
        t.Skip("integration disabled; set RUN_RECURPOST_INTEGRATION=1 to run")
    }

    base := strings.TrimSpace(os.Getenv("RECURPOST_BASE_URL"))
    email := strings.TrimSpace(os.Getenv("RECURPOST_EMAIL"))
    pass := strings.TrimSpace(os.Getenv("RECURPOST_PASSWORD"))
    accID := strings.TrimSpace(os.Getenv("RECURPOST_LINKEDIN_ACCOUNT_ID"))
    if base == "" || email == "" || pass == "" || accID == "" {
        t.Skipf("missing required envs; have base=%q email?=%t pass?=%t acc?=%t",
            base, email != "", pass != "", accID != "")
    }

    cli := rp.NewClient(rp.WithBaseURL(base))

    // Sanity check login works (and backend reachable)
    if _, err := cli.UserLogin.Login(context.Background(), rp.UserLoginRequest{EmailID: email, PassKey: pass}); err != nil {
        t.Fatalf("login failed: %v", err)
    }

    msg := strings.TrimSpace(os.Getenv("RECURPOST_LINKEDIN_TEST_MESSAGE"))
    if msg == "" {
        msg = "Ferret E2E LinkedIn test at " + time.Now().UTC().Format(time.RFC3339)
    }

    req := rp.PostContentRequest{
        EmailID:   email,
        PassKey:   pass,
        ID:        accID,
        Message:   msg,                      // fallback message
        LNMessage: "[LinkedIn] " + msg,      // LinkedIn-specific message
    }

    if mins := strings.TrimSpace(os.Getenv("RECURPOST_SCHEDULE_MINS")); mins != "" {
        // Some deployments accept schedule_date_time in e.g. "2006-01-02 15:04" format.
        // Use a common layout; adjust if your API expects a different one.
        when := time.Now().Add(5 * time.Minute).UTC().Format("2006-01-02 15:04")
        req.ScheduleDateTime = when
    }

    out, err := cli.Publishing.Post(context.Background(), req)
    if err != nil {
        t.Fatalf("post content failed: %v", err)
    }
    if out == nil || !out.Success || out.PostID == "" {
        t.Fatalf("unexpected response: %+v", out)
    }
}

func envBool(name string) bool {
    v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
    switch v {
    case "1", "true", "yes", "on":
        return true
    default:
        return false
    }
}
