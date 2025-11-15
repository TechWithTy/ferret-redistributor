package external

import (
    "context"
    "fmt"
    "os"
    "strings"

    ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
    "github.com/bitesinbyte/ferret/pkg/config"
)

type Instagram struct{}

func (m Instagram) Post(configData config.Config, post Post) error {
    // Build caption: title + hashtags (no clickable links in captions)
    caption := strings.TrimSpace(fmt.Sprintf("%s\n\n%s", post.Title, post.HashTags))
    // Try to resolve an OG image from the link
    imageURL, err := getOGImageURL(post.Link, configData)
    if err != nil || imageURL == "" {
        return err
    }
    cfg := ig.NewFromEnv()
    client := ig.New(cfg)
    ctx := context.Background()
    id, err := client.PostFeedImage(ctx, imageURL, caption)
    if err != nil { return err }
    // Optional: Post a first comment if provided via env
    if fc := strings.TrimSpace(os.Getenv("IG_FIRST_COMMENT_TEXT")); fc != "" {
        _ = client.PostFirstComment(ctx, id, fc)
    }
    return nil
}

// PostWithID posts to Instagram feed and returns the created media ID.
func (m Instagram) PostWithID(configData config.Config, post Post) (string, error) {
    caption := strings.TrimSpace(fmt.Sprintf("%s\n\n%s", post.Title, post.HashTags))
    imageURL, err := getOGImageURL(post.Link, configData)
    if err != nil || imageURL == "" { return "", err }
    cfg := ig.NewFromEnv()
    client := ig.New(cfg)
    ctx := context.Background()
    id, err := client.PostFeedImage(ctx, imageURL, caption)
    if err != nil { return "", err }
    if fc := strings.TrimSpace(os.Getenv("IG_FIRST_COMMENT_TEXT")); fc != "" {
        _ = client.PostFirstComment(ctx, id, fc)
    }
    return id, nil
}
