package generator

import (
    "encoding/json"
    "os"
)

// LoadVariants maps topic -> variants from JSON file (e.g., _data/variants.json).
func LoadVariants(path string) (map[string][]Variant, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    var m map[string][]Variant
    if err := json.Unmarshal(b, &m); err != nil { return nil, err }
    return m, nil
}

