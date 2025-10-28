Valkey Quick Start (no external deps)

Env
- `VALKEY_ADDR` (default `127.0.0.1:6379`)
- `VALKEY_PASSWORD` (optional)
- `VALKEY_DB` (optional integer)

Usage
```go
package main

import (
    "log"
    "time"
    "github.com/bitesinbyte/ferret/pkg/cache"
)

func main() {
    v, err := cache.NewValkey(cache.ValkeyConfig{})
    if err != nil { log.Fatal(err) }
    defer v.Close()
    if err := v.Ping(); err != nil { log.Fatal(err) }
    _ = v.Set("hello", "world", 10*time.Minute)
    val, ok, err := v.Get("hello")
    log.Println(val, ok, err)
}
```

Notes
- This client implements a tiny subset of RESP2 (PING/GET/SET/DEL, AUTH/SELECT).
- For production features (TLS, clustering), swap to `go-redis` and keep the same env keys.

