package factory

import (
    "strings"

    external "github.com/bitesinbyte/ferret/pkg/external"
)

// CreateSocialPoster returns a Poster implementation for the given platform name.
// Names are case-insensitive and support common aliases.
func CreateSocialPoster(name string) external.Poster {
    n := strings.ToLower(strings.TrimSpace(name))
    switch n {
    case "mastodon", "masto":
        return external.Mastodon{}
    case "twitter", "x":
        return external.Twitter{}
    case "linkedin", "li", "ln":
        return external.Linkedin{}
    case "facebook", "fb":
        return external.Facebook{}
    case "instagram", "ig":
        return external.Instagram{}
    default:
        // Fallback to Mastodon as a safe default
        return external.Mastodon{}
    }
}

