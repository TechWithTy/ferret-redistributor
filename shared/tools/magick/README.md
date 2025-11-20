## Magick Pipeline Tool

This folder contains helpers for running ImageMagick (MagickCore) pipelines backed by the `shared/models/magick_core.toon` schema. Agents can materialize a pipeline spec in JSON/TOON and invoke the script below to apply the operations with the `magick` CLI.

### Usage

```bash
python shared/tools/magick/pipeline.py \
  --spec /path/to/pipeline.json \
  --magick "/usr/local/bin/magick"   # optional
```

- `--spec` expects a JSON object shaped like `MagickPipeline` (see `shared/models/magick_core.toon`).
- `--magick` can override the binary path; otherwise the script uses `IMAGEMAGICK_BIN` or `PATH`.
- `--dry-run` prints the constructed command without executing.

Currently supported operations in the `operations` array:

| op_name      | Notes                                        |
|--------------|----------------------------------------------|
| `ResizeImage`| Uses `geometry.width/height` and optional filter |
| `CropImage`  | Requires geometry; supports offsets          |
| `ExtentImage`| Requires geometry; supports offsets          |
| `BlurImage`  | Uses `radius` / `sigma`                      |
| `SharpenImage` | Uses `radius` / `sigma`                   |
| `RotateImage`| Uses `angle`                                 |
| `FlipImage` / `FlopImage` | Mirrors vertically/horizontally |
| `ShearImage` | Uses geometry offsets                        |

Extend `pipeline.py` as needed to cover additional MagickCore operations and keep `shared/models/magick_core.toon` in sync.





