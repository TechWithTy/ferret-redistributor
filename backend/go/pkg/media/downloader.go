package media

import (
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type MediaDownloader struct {
    BasePath   string
    HTTPClient *http.Client
}

func New(basePath string) *MediaDownloader {
    return &MediaDownloader{BasePath: basePath, HTTPClient: &http.Client{Timeout: 30 * time.Second}}
}

// DownloadMedia downloads mediaURL to BasePath and returns the saved path.
func (m *MediaDownloader) DownloadMedia(mediaURL string) (string, error) {
    if m.BasePath == "" { m.BasePath = "downloads" }
    if err := os.MkdirAll(m.BasePath, 0o755); err != nil { return "", err }
    name := filenameFromURL(mediaURL)
    path := filepath.Join(m.BasePath, name)
    resp, err := m.HTTPClient.Get(mediaURL)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return "", err }
    f, err := os.Create(path)
    if err != nil { return "", err }
    defer f.Close()
    if _, err := io.Copy(f, resp.Body); err != nil { return "", err }
    return path, nil
}

func filenameFromURL(u string) string {
    u = strings.Split(u, "?")[0]
    parts := strings.Split(u, "/")
    if len(parts) == 0 || parts[len(parts)-1] == "" { return "media" }
    return parts[len(parts)-1]
}

