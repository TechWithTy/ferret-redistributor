package recurpost

import (
    "os"
    "testing"

    "github.com/joho/godotenv"
)

// TestMain loads environment variables from .env files before running tests,
// similar to a pytest conftest. It attempts both the repo root and CWD.
func TestMain(m *testing.M) {
    // Best-effort load from repo root when tests are invoked from project root
    _ = godotenv.Overload(".env")
    // Best-effort load when path is relative to this package location
    _ = godotenv.Overload("../../../../.env")

    os.Exit(m.Run())
}
