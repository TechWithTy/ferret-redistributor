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

## Model Files

### Notion Database Models
- `experiments_db.toon` - Experiments database schema
- `creative_assets.toon` - Creative Assets database schema
- `copy_calendar.toon` - Copy Calendar database schema
- `channels.toon` - Channels database schema
- `platforms.toon` - Platforms database schema
- `kpi_definitions.toon` - KPI Definitions database schema
- `kpi_progress.toon` - KPI Progress database schema
- `iterations_actions.toon` - Iterations/Actions database schema
- `scripts_variants.toon` - Scripts/Variants database schema

### Integration Models
- `postiz.toon` - Postiz MCP server data types for social media management
  - Post content, scheduling, campaigns, analytics, and connection management
  - Used by Postiz MCP prompts in `shared/prompt/experiments/sub_actions/postiz/`
  - See `shared/mcp/_docs/postiz-mcp.md` for API documentation

### Core Models
- `social_scale_core.toon` - Core social scale functionality models
- `magick_core.toon` - ImageMagick core operations
- `editly.toon` - Editly video editing models
- `media_edit_tracker.toon` - Media editing tracking

