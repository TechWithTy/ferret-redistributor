# Shared Tools

This directory hosts reusable tooling that multiple stacks (Go, Python, TS, C,
Rust) can share.

## Available Tools

| Script | Description |
| ------ | ----------- |
| `bootstrap_env.sh` | Creates a consistent dev shell (virtualenv + env vars) rooted at the repo so ad-hoc scripts can rely on predictable paths. |
| `run_codegen.py` | Entry point for future TOON → language generators. For now it validates `.toon` sources and prints guidance per target language. |
| `fetch_notion_dbs.py` | Downloads JSON snapshots for all Experiments-related Notion databases using env-configured tokens/IDs. |
| `cmd/imaging_ctools` | Thin wrapper around the `c/ImageMagick` submodule to run common image edits (resize, thumbnail, convert, optimize). |
| `cmd/prompt_stitch` | Compiles the master prompt with all sub-action prompts into a single artifact for MCP ingestion. |
| `video_editor.sh` | Proxy into `backend/typescript/tools` so you can run editly-based renders from the shared tooling surface. |
| `video/` | Editly pipeline helpers (`video/run_editly.py`) that wrap the backend TypeScript workspace. |
| `img/` | ImageMagick helpers (see `img/pipeline.py`) for executing `shared/models/magick_core.toon` specs. |
| `schemas/` | Auto-synced copies of all `.toon` files for MCP discovery (run `python shared/tools/sync_schemas.py`). |
| `Makefile` | Convenience targets (`make validate`, `make snapshot`) that wrap the Python helper and provide consistent status output in CI. |

## Usage

```bash
# 1. Setup shared env (creates .venv if missing)
shared/tools/bootstrap_env.sh

# 2. Validate TOON files / dry-run codegen
shared/tools/run_codegen.py --lang go

# 3. Run via make (used by CI)
make -C shared/tools validate

# 4. Pull the latest Notion data snapshots
export NOTION_API_TOKEN=secret_xxx
shared/tools/fetch_notion_dbs.py --output shared/generated/notion_snapshots

# 5. Run ImageMagick helpers
# (ensure `c/ImageMagick` is built or `magick` is installed locally)
go run ./shared/tools/cmd/imaging_ctools \
  -op resize \
  -input assets/source.png \
  -output assets/source-800w.jpg \
  -size 800x

# 5b. Stitch prompts
go run ./shared/tools/cmd/prompt_stitch \
  -master shared/prompt/experiments/master.poml \
  -subdir shared/prompt/experiments/sub_actions \
  -output shared/prompt/experiments/master_compiled.poml

# 6. Render video variants (choose shell or Python wrapper; both require pnpm + backend/typescript/tools deps)
shared/tools/video_editor.sh \
  --config backend/typescript/editly/examples/videos.json5 \
  --output /tmp/sample.mp4

# OR
python shared/tools/video/run_editly.py \
  --config backend/typescript/editly/examples/videos.json5 \
  --output /tmp/sample.mp4 \
  --fast

# 7. Run ImageMagick pipelines from TOON specs
python shared/tools/img/pipeline.py --spec /path/to/pipeline.json --dry-run

# 8. Refresh MCP schema copies
python shared/tools/sync_schemas.py
```

Feel free to extend these scripts with real code generators, additional
validators, or wrappers around external CLIs as we adopt them.

## Schema Resources

The TOON definitions live under `shared/models/` so MCP agents can reason about media workflows without scraping code. Run `python shared/tools/sync_schemas.py` to mirror them into `shared/tools/schemas/` for MCP servers that only scan the tools directory.

- `shared/models/editly.toon` → describes editly render specs (clips, layers, audio tracks, defaults). Agents populate this before calling `video_editor.sh` / `pnpm video:render`.
- `shared/models/magick_core.toon` → describes ImageMagick (MagickCore) pipelines (operations, geometry, filters, composite operators). Whenever you expose ImageMagick worklets via MCP, reference this schema so requests stay structured.

