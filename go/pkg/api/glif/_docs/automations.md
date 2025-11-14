# Glif Automation & Integration Notes

This document mirrors the automation briefs used by the other API SDKs. It
captures the constraints, required attribution, and CI guidance for the Glif
Simple API.

## Capabilities

- Run any Glif workflow via the Simple API (`https://simple-api.glif.app`) with
  positional or named inputs.
- Inspect workflow metadata, node names, and parameter schema from
  `https://glif.app/api/glifs?id=<workflowID>`.
- Toggle strict mode so workflows fail when required inputs are missing.

## Requirements

- Every UI that triggers a Glif workflow must show the “powered by Glif”
  attribution badge (provided in the partner brief).
- Simple API returns `200 OK` even on workflow failures. The SDK translates
  `response.error` into `ErrWorkflowFailed`; downstream consumers should handle
  that error explicitly.
- Rate limits are soft; Glif bills per run via the `price` field. Capture that in
  analytics to forecast spend.

## Environment Variables

| Variable         | Description                                 |
| ---------------- | ------------------------------------------- |
| `GLIF_API_TOKEN` | Bearer token from https://glif.app/settings |

Optional overrides:

- `GLIF_SIMPLE_API_BASE_URL`
- `GLIF_API_BASE_URL`

## GitHub Actions Snippet

```yaml
name: Glif Workflows
on:
  workflow_dispatch:
  schedule: [{ cron: "*/30 * * * *" }]

jobs:
  run:
    runs-on: ubuntu-latest
    env:
      GLIF_API_TOKEN: ${{ secrets.GLIF_API_TOKEN }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.21.x" }
      - run: go test ./go/pkg/api/glif/_tests
      - run: go run ./cmd/glif-runner
```

## Workflow Checklist

1. Fetch workflow metadata (`GetWorkflow`) to show node names in UI.
2. Validate user inputs client-side; provide either positional inputs OR named
   inputs, never both.
3. For strict workflows, set `RunWorkflowRequest.Strict = true` so Glif surfaces
   missing requirements.
4. Record `RunWorkflowResponse.Output` and `OutputFull`. The former is typically
   an asset URL; the latter contains provider-specific details for debugging.
5. Capture `RunWorkflowResponse.Price` for cost tracking.

