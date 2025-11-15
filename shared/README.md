# Shared Assets

The `shared/` tree centralizes artifacts that need to be reused across multiple
languages or runtimes:

- `shared/models/` — canonical TOON schemas that describe Social Scale entities.
  Downstream codegen (Go structs, Python dataclasses, TypeScript/Zod, SQL, etc.)
  should flow from these definitions.
- `shared/tools/` — scripts, CLIs, and helper configs that orchestrate model
  generation, validation, and external integrations.
- `shared/prompt/` — reusable POML prompt blueprints for MCP/TUI agents.
- `shared/resources/` — reference docs, specs, and knowledge base fragments
  that should be exposed read-only to orchestration agents.
- `shared/generated/` — machine-generated outputs (TypeScript, Go, SQL). This
  tree is **excluded** from MCP sharing to prevent accidental leakage of large
  artifacts or credentials baked into generated stubs.

Keeping these assets here prevents drift between Go, Python, TypeScript, Rust,
and C components while giving automation agents a single source of truth.

