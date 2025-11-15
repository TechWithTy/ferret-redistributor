# Go API SDKs

This directory hosts the API clients we vendor with Social Scale. Some are first-party packages that follow the same scaffolding as `postiz` or `glif`, while others are third-party SDKs pulled in as git submodules so we can track upstream changes directly.

## Structure

- `glif/` — native Go SDK for the Glif Simple API (beta) with typed helpers and tests.
- `fal/` — custom Fal.ai queue SDK (Call/Poll helpers + usage snippets).
- `rsshub/` — lightweight wrapper around RSSHub’s API and feed routes.
- `postiz/` — Social Scale’s generated Postiz Public API SDK (v1).
- `openroutergo/` — upstream OpenRouter SDK tracked as a git submodule from [`Deal-Scale/openroutergo`](https://github.com/Deal-Scale/openroutergo). Update via `git submodule update --remote go/pkg/api/openroutergo`.

## Submodules

After cloning the repository, initialize submodules (including OpenRouter) via:

```bash
git submodule update --init --recursive
```

When updating to the latest upstream SDK:

```bash
cd go/pkg/api/openroutergo
git fetch origin
git checkout main
git pull
```

Then commit the new submodule reference in the root repository.


