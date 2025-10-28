package linkedin_test

import (
    "io"
    "net/http"
    "strings"
    "testing"

    li "github.com/bitesinbyte/ferret/pkg/external/linkedin"
)

func mkResp(code int, body string) *http.Response {
    return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body))}
}

func TestMapHTTPError(t *testing.T) {
    if err := li.MapHTTPError(mkResp(401, "{}")); err != li.ErrUnauthorized { t.Fatal("401") }
    if err := li.MapHTTPError(mkResp(403, "{}")); err != li.ErrForbidden { t.Fatal("403") }
    if err := li.MapHTTPError(mkResp(429, "{}")); err != li.ErrRateLimited { t.Fatal("429") }
    if err := li.MapHTTPError(mkResp(500, "{}")); err != li.ErrServer { t.Fatal("500") }
    // 400 with message should return APIError
    err := li.MapHTTPError(mkResp(400, `{"message":"bad","serviceErrorCode":100}`))
    if _, ok := err.(li.APIError); !ok { t.Fatalf("expected APIError, got %T", err) }
}

