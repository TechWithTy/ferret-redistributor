# Shared TOON Models

Files in this directory define Social Scale's data contracts using TOON, our
language-agnostic schema format. Treat these as the **source of truth** for
entity shapes; downstream languages should rely on generated code rather than
hand-maintained structs.

## Workflow

1. Update the `.toon` definitions.
2. Run the forthcoming codegen tooling in `shared/tools/` to regenerate language-
   specific bindings (Go, Python, TypeScript, C, SQL, etc.).
3. Commit both the schema change and generated artifacts to keep all runtimes
   aligned.

> Until automated codegen is wired up, teams can reference the `.toon` file
> directly to keep implementations consistent.

