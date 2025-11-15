package storage

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type PublishedPost struct {
    Platform    string    `json:"platform"`
    ID          string    `json:"id"`
    Link        string    `json:"link,omitempty"`
    ContentType string    `json:"content_type,omitempty"`
    PublishedAt time.Time `json:"published_at"`
}

// AppendPublishedPost appends a JSON line to data/published_posts.jsonl under repo root.
func AppendPublishedPost(rec PublishedPost) error {
    dir := filepath.Join("data")
    if err := os.MkdirAll(dir, 0o755); err != nil { return err }
    path := filepath.Join(dir, "published_posts.jsonl")
    f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
    if err != nil { return err }
    defer f.Close()
    b, err := json.Marshal(rec)
    if err != nil { return err }
    if _, err := f.Write(append(b, '\n')); err != nil { return err }
    return nil
}

