package behiiv_test

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    bh "github.com/bitesinbyte/ferret/pkg/external/behiiv"
)

func TestGetAggregateStats(t *testing.T) {
    pub := "pub_00000000-0000-0000-0000-000000000000"
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") != "Bearer tok" { t.Fatalf("missing auth header") }
        if r.URL.Path != "/v2/publications/"+pub+"/posts/aggregate_stats" {
            t.Fatalf("path: %s", r.URL.Path)
        }
        _ = json.NewEncoder(w).Encode(map[string]any{
            "data": map[string]any{
                "stats": map[string]any{
                    "email": map[string]any{"recipients": 100, "opens": 50},
                    "web": map[string]any{"views": 200},
                    "clicks": []any{},
                },
            },
        })
    }))
    defer ts.Close()

    cfg := bh.Config{BaseURL: ts.URL, Version: "v2", Token: "tok"}
    c := bh.New(cfg)
    resp, err := c.GetAggregateStats(context.Background(), pub)
    if err != nil { t.Fatalf("err: %v", err) }
    if resp.Data.Stats.Email.Recipients != 100 { t.Fatalf("recipients: %d", resp.Data.Stats.Email.Recipients) }
    if resp.Data.Stats.Web.Views != 200 { t.Fatalf("views: %d", resp.Data.Stats.Web.Views) }
}

func TestListPosts(t *testing.T) {
    pub := "pub_1"
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/v2/publications/"+pub+"/posts" { t.Fatalf("path: %s", r.URL.Path) }
        _ = json.NewEncoder(w).Encode(map[string]any{
            "data": []any{
                map[string]any{
                    "id": "post_1",
                    "title": "T",
                    "subtitle": "S",
                    "authors": []any{"Clark Kent"},
                    "created": 1666800076,
                    "status": "draft",
                    "split_tested": true,
                    "subject_line": "Check this out",
                    "preview_text": "More news on the horizon",
                    "slug": "slug",
                    "thumbnail_url": "thumbnail_url",
                    "web_url": "web_url",
                    "audience": "free",
                    "platform": "web",
                    "content_tags": []any{"content_tags"},
                    "hidden_from_feed": true,
                    "publish_date": 1666800076,
                    "displayed_date": 1666800076,
                    "meta_default_description": "A post with great content",
                    "meta_default_title": "My great post",
                    "content": map[string]any{
                        "free": map[string]any{"web": "<html>", "email": "<html>", "rss": "<xml>"},
                        "premium": map[string]any{"web": "<html>", "email": "<html>"},
                    },
                    "stats": map[string]any{
                        "email": map[string]any{"recipients": 100, "opens": 50, "unique_opens": 45},
                        "web": map[string]any{"views": 200, "clicks": 40},
                        "clicks": []any{ map[string]any{ "url": "https://www.google.com", "email": map[string]any{"clicks":10, "unique_clicks":8, "click_through_rate":80}, "web": map[string]any{"clicks":40, "unique_clicks":40, "click_through_rate":20}, "total_clicks":50, "total_unique_clicks":48, "total_click_through_rate":40 } },
                    },
                },
            },
            "limit": 1, "page": 1, "total_results": 1, "total_pages": 1,
            "links": map[string]any{"next": "", "prev": ""},
        })
    }))
    defer ts.Close()
    cfg := bh.Config{BaseURL: ts.URL, Version: "v2", Token: "tok"}
    c := bh.New(cfg)
    resp, err := c.ListPosts(context.Background(), pub, 0, 0)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(resp.Data) != 1 || resp.Data[0].ID != "post_1" { t.Fatalf("resp: %+v", resp) }
    if resp.Data[0].Stats.Email.Recipients != 100 { t.Fatalf("stats parse failed: %+v", resp.Data[0].Stats) }
    if resp.Limit != 1 || resp.TotalResults != 1 { t.Fatalf("meta parse failed: %+v", resp) }
}
