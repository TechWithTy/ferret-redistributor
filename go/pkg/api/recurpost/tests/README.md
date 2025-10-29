RecurPost Client Tests

This test suite exercises the RecurPost API client against a local mock server. It NEVER calls the real RecurPost API or creates any posts.

Environment

- Define credentials in the project root `.env` file or the shell before running:
  - `RECURPOST_EMAIL`
  - `RECURPOST_PASSWORD`
- The tests auto-load `../../../../.env` so you can store the variables in `go/.env`.

Run

- Run all tests for this suite only:
  - `go test -v ./pkg/api/recurpost/tests`
- Or with the parent package (which has no tests):
  - `go test -v ./pkg/api/recurpost/...`

Notes

- Tests skip if the env vars are missing.
- All endpoints are mocked with `httptest.Server` to avoid network calls and side effects.

