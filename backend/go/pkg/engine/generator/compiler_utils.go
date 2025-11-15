package generator

import (
    "fmt"
    "os"
    "path/filepath"
)

func ensureDir(p string) error {
    dir := p
    if ext := filepath.Ext(p); ext != "" {
        dir = filepath.Dir(p)
    }
    if dir == "." || dir == "" {
        return nil
    }
    return os.MkdirAll(dir, 0o755)
}

func writeFile(outPath string, data []byte) error {
    if err := ensureDir(outPath); err != nil {
        return fmt.Errorf("ensureDir failed: %w", err)
    }
    return os.WriteFile(outPath, data, 0o644)
}

