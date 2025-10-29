package cache

import "fmt"

type RedisCache struct{}

func (r RedisCache) Save(key, value string) {
    fmt.Printf("CacheEngine: Saved [%s] -> %s\n", key, value)
}

func (r RedisCache) Get(key string) string {
    fmt.Printf("CacheEngine: Retrieved key [%s]\n", key)
    return "mock-value"
}

