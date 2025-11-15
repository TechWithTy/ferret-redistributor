## TypeScript Video Tools

This package wraps our local `typescript/editly` workspace to provide a single entry point for creating or iterating on video experiments.

### Setup

```bash
cd typescript/tools
pnpm install
```

The dependency `editly` points to `../editly`, so make sure that package builds successfully (see its README for FFmpeg requirements).

### Render a Video

```bash
pnpm video:render -- --config ../editly/examples/videos.json5 --output /tmp/sample.mp4
```

Available flags:

| Flag | Description |
| ---- | ----------- |
| `--config, -c` | Path to an editly JSON / JSON5 config file. Required. |
| `--output, -o` | Override `outPath` in the config. |
| `--width` / `--height` | Override render dimensions. |
| `--fast` | Enables editly’s `fast` mode (skips some transitions). |

### Usage in Workflows

- Reference this tool from prompts/sub-actions when creating A/B video variants.
- Consider checking in canonical config templates next to experiment documentation for reproducibility.



