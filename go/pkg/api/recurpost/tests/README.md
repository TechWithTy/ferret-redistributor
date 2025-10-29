RecurPost Client Tests

This test suite exercises the RecurPost API client against a local mock server. It NEVER calls the real RecurPost API or creates any posts.

Environment

- Define credentials in the project root `.env` file or the shell before running:
  - `RECURPOST_EMAIL`
  - `RECURPOST_PASSWORD`
- Integration (live) test requires:
  - `RUN_RECURPOST_INTEGRATION=1`
  - `RECURPOST_BASE_URL` (e.g., https://api.recurpost.com)
  - `RECURPOST_LINKEDIN_ACCOUNT_ID`
- Tests auto-load `.env` via TestMain (see `test_setup.go`). You can put the variables in `go/.env`.

Run

- Run all tests for this suite only:
  - `go test -v ./pkg/api/recurpost/tests`
- Or with the parent package (which has no tests):
  - `go test -v ./pkg/api/recurpost/...`
- Run LinkedIn unit test only:
  - `go test -v ./pkg/api/recurpost/tests -run ^TestLinkedinPost$`
- Run LinkedIn LIVE test only (uses env/.env):
  - `go test -v -count=1 ./pkg/api/recurpost/tests -run ^TestLinkedinPost_Live$`

Notes

- Tests skip if the env vars are missing.
- All endpoints are mocked with `httptest.Server` to avoid network calls and side effects.
