package instagram

import (
    "context"
    "net/http"
    "strings"
    "testing"
)

func TestBuildCreateMediaForm(t *testing.T) {
    c := &Client{cfg: Config{AccessToken: "tok"}}
    f := c.buildCreateMediaForm(createMediaParams{
        Caption: "cap",
        ImageURL: "http://img",
        ThumbOffsetSeconds: 5,
        DisableComments: true,
        CoverURL: "http://cover",
        IsCarousel: true,
        Children: []string{"1","2"},
    })
    if f.Get("access_token") != "tok" { t.Fatal("missing token") }
    if f.Get("caption") != "cap" { t.Fatal("missing caption") }
    if f.Get("image_url") != "http://img" { t.Fatal("missing image_url") }
    if f.Get("thumb_offset") != "5" { t.Fatal("thumb_offset not set") }
    if f.Get("disable_comments") != "true" { t.Fatal("disable_comments not set") }
    if f.Get("cover_url") != "http://cover" { t.Fatal("cover_url not set") }
    if f.Get("media_type") != "CAROUSEL" { t.Fatal("media_type not carousel") }
    if got := f.Get("children"); got != "1,2" { t.Fatalf("children got %s", got) }
}

func TestMapHTTPError(t *testing.T) {
    r := &http.Response{StatusCode: 400, Body: ioNopCloser(`{"error":{"message":"bad","code":100}}`)}
    if err := mapHTTPError(r); err == nil { t.Fatal("expected error") }
}

type nopCloser struct{ *strings.Reader }
func (n nopCloser) Close() error { return nil }
func ioNopCloser(s string) nopCloser { return nopCloser{strings.NewReader(s)} }

func TestPostFirstCommentValidation(t *testing.T) {
    c := New(Config{})
    if err := c.PostFirstComment(context.Background(), "", "msg"); err == nil {
        t.Fatal("expected validation error")
    }
}

