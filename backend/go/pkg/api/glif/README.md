# Glif Simple API SDK (Go)

Typed helper for the Glif Simple API and the public REST API (`/api/glifs`). This package mirrors the structure of the existing SDKs under `pkg/api/*` and documents the beta-specific behaviors and exceptions described in the Glif API reference.

## API Surface

- Simple API base URL: `https://simple-api.glif.app`
- REST API base URL: `https://glif.app/api`
- Auth: Bearer token (`Authorization: Bearer <token>`)
- Status codes: Simple API always returns **200 OK**, even when the workflow fails. Inspect the `error` field in the JSON body.
- Strict mode: append `?strict=1` (set `RunWorkflowRequest.Strict=true`) to force failures when required inputs are missing.
- Visibility: defaults to `PRIVATE`; set `RunWorkflowRequest.Visibility` to `PUBLIC` if your integration explicitly needs public runs.

## Usage

```go
ctx := context.Background()

cli := glif.NewClient(
  os.Getenv("GLIF_API_TOKEN"),
  glif.WithSimpleAPIBaseURL("https://simple-api.glif.app"),
)

resp, err := cli.RunWorkflow(ctx, glif.RunWorkflowRequest{
  WorkflowID: "clgh1vxtu0011mo081dplq3xs",
  Inputs:     []string{"cute friendly oval shaped bot friend"},
})
if err != nil {
  log.Fatal(err)
}
fmt.Println("Output URL:", resp.Output)
```

### Named inputs

```go
resp, err := cli.RunWorkflow(ctx, glif.RunWorkflowRequest{
  WorkflowID: "clkbasluf0000mi08h541a3j4",
  NamedInputs: map[string]string{
    "username": "jamiedubs",
  },
  Strict: true,
})
```

### Fetch workflow graph / node names

```go
wf, err := cli.GetWorkflow(ctx, "clkbasluf0000mi08h541a3j4")
if err != nil {
  log.Fatal(err)
}
for _, node := range wf.Data.Nodes {
  fmt.Println(node.Name, node.Type)
}
```

## Exceptions & Beta Constraints

- **200-on-error:** the Simple API currently never sets non-2xx codes; all workflow-level errors are surfaced in `RunWorkflowResponse.Error`. The SDK wraps these as `ErrWorkflowFailed`.
- **Conflicting inputs:** the API accepts either positional or named inputs. The SDK enforces this and returns `ErrConflictingInputModes` if both are supplied.
- **Power attribution:** Glif requires "powered by glif.app" + logo in your UI. Ensure your product surface includes the provided badge asset (see docs link in the product brief).
- **Rate limits:** SDK does not throttle. Monitor `price` and API usage; contact Glif if you need higher limits.
- **CORS:** only the Simple API allows browser requests. Other endpoints should be proxied through your own backend if needed.

## Testing

Integration-style tests now live under `_tests` to match the rest of the Social Scale SDKs. They use `httptest.Server` and require no real Glif credentials:

```bash
cd go
go test ./pkg/api/glif/_tests
```
