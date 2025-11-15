# Shared Resources

Drop any reusable reference material here (knowledge base notes, JSON fixtures,
workflow specs, etc.). Everything in this directory is automatically exposed to
the MCP shared-content server as read-only resources, so avoid placing secrets
or credentials in this tree. Use the `_mcp_deny`, `_auth`, or `_service`
suffixes if you need to hide a file from the server without deleting it.

## Prompt Stitching (Experiments)

Use `shared/tools/cmd/prompt_stitch` to compile the top-level Experiments prompt and every sub-action prompt into a single artifact (default: `shared/prompt/experiments/master_compiled.poml`).

```bash
# Run locally
go run ./shared/tools/cmd/prompt_stitch \
  -master shared/prompt/experiments/master.poml \
  -subdir shared/prompt/experiments/sub_actions \
  -output shared/prompt/experiments/master_compiled.poml
```

### CI/CD Integration

Add a lightweight job (example: GitHub Actions) so stitched prompts remain current on every pull request:

```yaml
jobs:
  prompt-stitch:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go run ./shared/tools/cmd/prompt_stitch
      - name: Verify diff
        run: |
          if ! git diff --quiet -- shared/prompt/experiments/master_compiled.poml; then
            echo "::error::Stitched prompt is out of date. Run go run ./shared/tools/cmd/prompt_stitch locally."
            exit 1
          fi
```

In other CI providers, replicate the same steps: install Go, run the stitcher, and fail the build if the generated file changes.

## Image Editing via c/ImageMagick

The repository vendors ImageMagick as a submodule under `c/ImageMagick`. Use `shared/tools/cmd/imaging_ctools` to run common edits without having to author CLI commands manually.

### Build / Dependencies

```bash
cd c/ImageMagick
./configure --disable-dependency-tracking
make -j$(nproc)
export IMAGEMAGICK_BIN="$PWD/utilities/magick"
```

You can also install ImageMagick system-wide; the helper falls back to `magick` on `$PATH`.

### Usage

```bash
go run ./shared/tools/cmd/imaging_ctools \
  -op thumbnail \
  -input docs/source.png \
  -output docs/source-thumb.jpg \
  -size 400x400
```

Supported ops: `resize`, `thumbnail`, `convert`, `optimize` (with `-quality`).

