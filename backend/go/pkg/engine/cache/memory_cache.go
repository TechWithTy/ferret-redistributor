package cache

import "fmt"

type MemoryCache struct{ store map[string]string }

func (m *MemoryCache) Save(key, value string) {
    if m.store == nil {
        m.store = map[string]string{}
    }
    m.store[key] = value
    fmt.Printf("CacheEngine: [mem] Saved [%s]\n", key)
}

func (m *MemoryCache) Get(key string) string {
    if m.store == nil {
        return ""
    }
    return m.store[key]
}

