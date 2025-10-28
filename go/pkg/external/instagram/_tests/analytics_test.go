package instagram_test

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
)

func TestGetMediaBasicAndInsights(t *testing.T) {
    mediaID := "1789012345"
    igUser := "1234"
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch {
        case r.Method == http.MethodGet && r.URL.Path == "/"+mediaID:
            _ = json.NewEncoder(w).Encode(map[string]any{
                "id": mediaID, "like_count": 10, "comments_count": 3, "play_count": 100, "save_count": 2,
            })
        case r.Method == http.MethodGet && r.URL.Path == "/"+mediaID+"/insights":
            _ = json.NewEncoder(w).Encode(map[string]any{
                "data": []map[string]any{
                    {"name": "impressions", "period": "lifetime", "values": []map[string]any{{"value": 123}}},
                    {"name": "reach", "period": "lifetime", "values": []map[string]any{{"value": 77}}},
                },
            })
        case r.Method == http.MethodGet && r.URL.Path == "/"+igUser+"/insights":
            _ = json.NewEncoder(w).Encode(map[string]any{
                "data": []map[string]any{
                    {"name": "impressions", "values": []map[string]any{{"value": 321.0, "end_time": "2024-01-01T00:00:00Z"}}},
                },
            })
        default:
            http.NotFound(w, r)
        }
    }))
    defer ts.Close()

    cfg := ig.Config{BaseURL: ts.URL, AccessToken: "tok", IGUserID: igUser}
    client := ig.New(cfg)
    ctx := context.Background()

    basic, err := client.GetMediaBasic(ctx, mediaID)
    if err != nil { t.Fatalf("basic err: %v", err) }
    if basic.LikeCount != 10 || basic.CommentsCount != 3 { t.Fatalf("unexpected basic: %+v", basic) }

    ins, err := client.GetMediaInsights(ctx, mediaID, []string{"impressions","reach"})
    if err != nil { t.Fatalf("insights err: %v", err) }
    if len(ins) != 2 { t.Fatalf("insights len: %d", len(ins)) }

    uins, err := client.GetUserInsights(ctx, []string{"impressions"}, "day", "", "")
    if err != nil { t.Fatalf("user insights err: %v", err) }
    if len(uins) == 0 || uins[0].Value != 321 { t.Fatalf("unexpected user insights: %+v", uins) }
}

