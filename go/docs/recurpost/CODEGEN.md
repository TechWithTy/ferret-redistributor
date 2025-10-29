# RecurPost SDK Codegen

This repo includes an OpenAPI spec for the RecurPost API at:
- `pkg/api/recurpost/openapi_new.yaml`

Generate Go types/clients with oapi-codegen:

1) Install (one-time)

```
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

2) Generate

```
# Types + client in the recurpost package (adjust output path as desired)
oapi-codegen -generate types,client \
  -package recurpost \
  -o pkg/api/recurpost/generated_oapi.go \
  pkg/api/recurpost/openapi_new.yaml
```

Notes
- Prefer committing the spec and generated code together to reduce drift.
- You can also generate server stubs (`-generate chi-server` or `gin`) if needed.
- Keep hand-written files focused on helpers, not duplicating generated models.

