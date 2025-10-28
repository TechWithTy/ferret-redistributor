package external

import "github.com/bitesinbyte/ferret/pkg/config"

type Post struct {
	Title       string
	Link        string
	Description string
	HashTags    string
}
type Poster interface {
    Post(configData config.Config, post Post) error
}

// PosterWithID is an optional extension that returns a created post ID/URN
// to enable analytics fetching.
type PosterWithID interface {
    PostWithID(configData config.Config, post Post) (string, error)
}
