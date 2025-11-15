## TOON Support Matrix

This repository now vendors the canonical language implementations of the
Token-Oriented Object Notation (TOON) format so workflow payloads can be
serialized efficiently across stacks.

| Language      | Location                         | Upstream Repo |
| ------------- | -------------------------------- | ------------- |
| Go            | `go/pkg/external/gotoon`         | https://github.com/alpkeskin/gotoon |
| Python        | `python/toon_support`            | Local helpers (see `encoder.py`) |
| TypeScript    | `backend/typescript/editly` (placeholder) | Bring-your-own until a dedicated TOON TS lib exists |
| C / C++       | `c/ctoon` and `c/ImageMagick`    | https://github.com/mohammadraziei/ctoon |
| Rust          | `rust/toon-rust`                 | https://github.com/toon-format/toon-rust |

> **Note:** Until the repository is moved, the Go module path still uses `github.com/bitesinbyte/ferret/...`. Imports will be updated once Social Scale’s canonical path is finalized.

### Python Usage

```python
from python.toon_support import encode_to_toon

payload = {"workflow": "social-post", "posts": [{"id": 1, "status": "queued"}]}
print(encode_to_toon("social-scale", payload))
```

### Go Usage

```go
import "github.com/bitesinbyte/ferret/go/pkg/external/gotoon"

encoded, err := gotoon.Encode(payload)
```

### Rust / C

Rust and C clients can build against the vendored submodules. Run `git submodule update --init --recursive` after cloning to hydrate the sources.

