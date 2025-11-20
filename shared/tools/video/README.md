## Video Pipeline Tool

`run_editly.py` bridges MCP tooling with the backend TypeScript editly workspace (`backend/typescript/tools`). It accepts an editly JSON/JSON5 config (matching `shared/models/editly.toon`) and invokes the existing `pnpm video:render` script with optional overrides.

### Requirements

- Node.js + pnpm installed (Corepack works too).
- `backend/typescript/tools` dependencies installed (`pnpm install`), plus the local `editly` workspace.

### Usage

```bash
python shared/tools/video/run_editly.py \
  --config backend/typescript/editly/examples/videos.json5 \
  --output /tmp/sample.mp4 \
  --fast \
  --width 720 \
  --height 1280
```

Flags map directly to editly CLI options:

| Flag | Description |
| ---- | ----------- |
| `--config, -c` | Path to an editly JSON/JSON5 spec. Required. |
| `--output, -o` | Overrides `outPath`. |
| `--width/--height` | Override render dimensions. |
| `--fast` | Enables editly fast/preview mode. |
| `--pnpm` | Explicit path to pnpm (defaults to `PNPM_BIN` or PATH). |
| `--cwd` | Working directory (defaults to `backend/typescript/tools`). |
| `--dry-run` | Print the pnpm command without executing it. |

Agents should populate the `editly.toon` schema, write it to disk (JSON/JSON5), then call this script via MCP tooling.





