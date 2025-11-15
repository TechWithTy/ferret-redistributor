# Engine-Oriented Workflow Runtime (EOWR)

This guide shows how to author EOWR pseudo workflows, compile them to runnable Go, and export n8n-compatible JSON. It covers basic and advanced patterns so you can go from a single .pseudo file to code and a visual workflow in minutes.

## Quick Start

- Place pseudo files in `pkg/engine/workflows`.
- Use the CLI to compile and export.
- Import the generated JSON into n8n.

Examples below assume Windows shell. On macOS/Linux, drop `.exe` and adjust paths as needed.

### 1) Author a workflow

Create `pkg/engine/workflows/lead_import.pseudo`:

```
workflow "Lead Import" version 1.0:
  trigger "Webhook" at "/new-lead"
  action "Transform Lead" using "workers/ai_worker"
  action "Push to CRM" using "queue/nats_engine"
  on_error "Notify" using "telemetry/sentry_engine"
  connect "Webhook" -> "Transform Lead"
  connect "Transform Lead" -> "Push to CRM"
  connect "Push to CRM" -> "Notify"
```

### 2) Compile to Go and export to n8n JSON

- Compile all pseudo => Go:
  `go run ./cmd/eowr compile-all`

- Export all pseudo => JSON:
  `go run ./cmd/eowr export-all`

Outputs land in:
- Go: `pkg/engine/workflows/compiled/<name>.go`
- JSON: `pkg/engine/workflows/exports/<name>.json`

You can also operate on a single file:
- `go run ./cmd/eowr compile pkg/engine/workflows/lead_import.pseudo pkg/engine/workflows/compiled/lead_import.go`
- `go run ./cmd/eowr export  pkg/engine/workflows/lead_import.pseudo pkg/engine/workflows/exports/lead_import.json`

Makefile shortcuts:
- `make eowr-compile IN=pkg/engine/workflows/lead_import.pseudo OUT=pkg/engine/workflows/compiled/lead_import.go`
- `make eowr-export  IN=pkg/engine/workflows/lead_import.pseudo OUT=pkg/engine/workflows/exports/lead_import.json`
- `make eowr-bootstrap` (compile-all + export-all)

### 3) Import into n8n

In your n8n instance:
- Create workflow → Import from file
- Select the generated JSON (e.g., `pkg/engine/workflows/exports/lead_import.json`).
- Review connections and node parameters; fill credentials where required.

## Pseudo DSL Reference

Supported node types and statements (one per line):

- `workflow "Name" version X.Y:`
  - Declares the workflow name and version.

- `trigger "Webhook" at "/path"`
  - Entry node. `at` is optional; used for webhook path or schedule name.

- `action "Operation" using "package/ref"`
  - Operation node. `using` refers to an engine or worker path (informational in JSON export; used for hints in Go).

- `switch "Label"`
  - Conditional branch entry (rules attached via n8n after import). Basic exporter maps type only.

- `merge "Name"`
  - Combine branches back together.

- `subflow "Other Flow" using "path/to/flow"`
  - Nested workflow reference.

- `on_error "Handler" using "telemetry/sentry_engine"`
  - Error handler node.

- `connect "A" -> "B"`
  - Connect nodes by name.

- `connect "A" -> ["B", "C"]`
  - Fan-out connections to multiple targets.

Notes:
- Names must match exactly between node declarations and `connect` statements.
- You can declare nodes in any order; connections define the graph.

## Basic Example

```
workflow "ETL" version 1.0:
  trigger "Cron"
  action "Extract" using "workers/data_worker"
  action "Transform" using "workers/data_worker"
  action "Load" using "queue/pulsar_engine"
  connect "Cron" -> "Extract"
  connect "Extract" -> "Transform"
  connect "Transform" -> "Load"
```

Compile/export:
- `go run ./cmd/eowr compile-all`
- `go run ./cmd/eowr export-all`

## Advanced Example (Branching, Merge, Subflow, Error)

```
workflow "Lead Import and Enrichment" version 3.2:
  trigger "Webhook" at "/new-lead"
  action "Parse JSON" using "core/actions/parse"
  action "Enrich Data" using "workers/ai_worker"
  switch "Check Score"
  action "Push to CRM" using "queue/nats_engine"
  action "Send Alert" using "telemetry/sentry_engine"
  merge "End"
  subflow "CRM Sync" using "workflows/crm_sync"
  on_error "Alert" using "telemetry/sentry_engine"

  connect "Webhook" -> "Parse JSON"
  connect "Parse JSON" -> "Enrich Data"
  connect "Enrich Data" -> "Check Score"
  connect "Check Score" -> ["Push to CRM", "Send Alert"]
  connect ["Push to CRM", "Send Alert"] -> "End"
  connect "End" -> "CRM Sync"
  connect "End" -> "Alert"
```

Exporter behavior:
- Nodes are mapped to n8n types: trigger → webhook, action → httpRequest/function, switch → switch, merge → merge, subflow → executeWorkflow, on_error → code.
- Rules/parameters beyond name/type may need to be set inside n8n after import (e.g., switch rules, credentials).

## Tips & Conventions

- File placement: keep `.pseudo` files in `pkg/engine/workflows` so `compile-all` and `export-all` pick them up automatically.
- Naming: keep node names unique within a workflow to avoid connection ambiguity.
- Error handling: include an `on_error` node and connect to it from terminal steps if you want error routes visible in n8n.
- Subflows: reference other workflow names in `subflow` and connect to them as needed. The exporter emits the correct node type; you’ll select the target workflow in n8n UI.
- Credentials: exporters don’t embed secrets; set API creds inside n8n after import.

## Troubleshooting

- CLI not found: use `go run ./cmd/eowr ...` from repo root.
- No outputs generated: ensure your `.pseudo` file matches the syntax above and lives under `pkg/engine/workflows`.
- Broken imports in generated Go: the translator uses module path `github.com/bitesinbyte/ferret`. Check `go.mod` and the engine package paths if you moved files.
- n8n connections look empty: verify node names in your `connect` lines match the declared node names exactly.

## Internals

- Parser: `pkg/engine/generator/pseudo_parser.go` extracts workflow name, nodes, and connections.
- Go generator: `pkg/engine/generator/go_translator.go` emits a runnable stub using engines and workers.
- n8n exporter: `pkg/engine/generator/n8n_exporter.go` builds minimal n8n JSON (nodes + connections). Some parameters (e.g., switch rules) are best edited in the n8n UI post-import.

