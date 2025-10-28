package cache

import (
    "encoding/json"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// FileCache is a simple JSON-on-disk cache with TTL per entry.
// It writes files under baseDir (default: _data_cache). Entries are stored as
// {"expires_at": unixSeconds, "value": <arbitrary JSON>}.
type FileCache struct {
    baseDir string
}

// NewFileCache creates a cache rooted at baseDir, creating it if needed.
func NewFileCache(baseDir string) (*FileCache, error) {
    if baseDir == "" {
        baseDir = "_data_cache"
    }
    if err := os.MkdirAll(baseDir, 0o755); err != nil {
        return nil, err
    }
    return &FileCache{baseDir: baseDir}, nil
}

// GetBytes returns the raw JSON payload for key if present and not expired.
func (c *FileCache) GetBytes(key string) ([]byte, bool) {
    p := c.pathFor(key)
    b, err := os.ReadFile(p)
    if err != nil {
        return nil, false
    }
    var wrap struct{
        ExpiresAt float64       `json:"expires_at"`
        Value     json.RawMessage `json:"value"`
    }
    if err := json.Unmarshal(b, &wrap); err != nil {
        _ = os.Remove(p)
        return nil, false
    }
    if wrap.ExpiresAt > 0 && time.Now().After(time.Unix(int64(wrap.ExpiresAt), 0)) {
        _ = os.Remove(p)
        return nil, false
    }
    return []byte(wrap.Value), true
}

// GetJSON unmarshals the cached JSON into dst if present and valid.
func (c *FileCache) GetJSON(key string, dst any) (bool, error) {
    raw, ok := c.GetBytes(key)
    if !ok {
        return false, nil
    }
    if err := json.Unmarshal(raw, dst); err != nil {
        return false, err
    }
    return true, nil
}

// SetBytes stores raw JSON bytes with a TTL.
func (c *FileCache) SetBytes(key string, raw []byte, ttl time.Duration) error {
    if ttl <= 0 {
        return errors.New("ttl must be > 0")
    }
    p := c.pathFor(key)
    tmp := p + ".tmp"
    wrap := struct{
        ExpiresAt float64 `json:"expires_at"`
        Value     json.RawMessage `json:"value"`
    }{ExpiresAt: float64(time.Now().Add(ttl).Unix()), Value: json.RawMessage(raw)}
    b, err := json.Marshal(wrap)
    if err != nil {
        return err
    }
    if err := os.WriteFile(tmp, b, 0o644); err != nil {
        return err
    }
    return os.Rename(tmp, p)
}

// SetJSON marshals v to JSON and stores it.
func (c *FileCache) SetJSON(key string, v any, ttl time.Duration) error {
    b, err := json.Marshal(v)
    if err != nil {
        return err
    }
    return c.SetBytes(key, b, ttl)
}

// Delete removes a key if present.
func (c *FileCache) Delete(key string) error {
    p := c.pathFor(key)
    if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
        return err
    }
    return nil
}

func (c *FileCache) pathFor(key string) string {
    safe := strings.NewReplacer("/", "_", "\\", "_", ":", "_").Replace(key)
    return filepath.Join(c.baseDir, safe+".json")
}

