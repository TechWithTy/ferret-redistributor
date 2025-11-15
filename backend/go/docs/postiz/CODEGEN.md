# Postiz SDK Codegen

This repo includes an OpenAPI spec for the Postiz Public API at:
- `pkg/api/postiz/openapi.yaml`

Generate Go types/clients with oapi-codegen:

1) Install (one-time)

```
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

2) Generate

```
# Types + client in the postiz package (adjust output path as desired)
oapi-codegen -generate types,client \
  -package postiz \
  -o pkg/api/postiz/generated_oapi.go \
  pkg/api/postiz/openapi.yaml
```

Notes
- Prefer committing the spec and generated code together to reduce drift.
- Keep hand-written files focused on helpers, not duplicating generated models.

