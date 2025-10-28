# CyberOni Google Trends SDK (Alpha Scaffold)

This package is an alpha scaffold for a future Google Trends SDK. At the time of writing, there is no official public Google Trends API generally available. Many projects use the community `pytrends` package or a proxy service. This SDK defines a clean, typed interface and utilities so you can later plug in an implementation (e.g., `pytrends`) or an official API when it becomes available.

- Folder: `backend/app/core/landing_page/google_trends/`
- Client: `client.py` exposes `GoogleTrendsClient`
- Status: Interface only; methods currently raise `NotImplementedError`

## Alpha Status and Philosophy
- API endpoints, auth model, and schemas are not finalized. The SDK is designed to be adaptable.
- Config is driven by environment variables with sensible defaults; see `.env.example`.
- Utilities mirror patterns used in GA4/GSC SDKs for consistency.

## Install (local dev)
No extra runtime deps are required for the scaffold. If you wire a backend like `pytrends`, add it to `pyproject.toml` or your environment.

## Environment Variables
See `backend/app/core/landing_page/google_trends/.env.example` for all options. Key variables:

- `TRENDS_API_BASE_URL` (default: `https://trends.googleapis.com/v1/`)
- `TRENDS_API_KEY` (optional; if an API key model is used)
- `TRENDS_SCOPE` (optional; if OAuth scopes are used)
- `TRENDS_HTTP_TIMEOUT` (default: `30` seconds)
- `TRENDS_USER_AGENT` (default: `CyberOni-GTrends-SDK/0.1`)
- Feature flags: `TRENDS_ALPHA_MODE`, `TRENDS_ENABLE_MERGE`, `TRENDS_DEFAULT_INTERVAL`

## Quickstarts
Two quickstarts demonstrate request building and intended usage. Methods are stubs and will print the built request until an implementation is provided.

- `scripts/quickstart_trends_alpha.py` — single-term daily interest for the last 30 days.
- `scripts/quickstart_trends_compare.py` — compare multiple terms.

Run (from repo root or adjust PYTHONPATH):
```
# Example (Windows PowerShell)
$env:TRENDS_TERM="dealscale"; 
python backend/app/core/landing_page/google_trends/scripts/quickstart_trends_alpha.py

$env:TRENDS_TERMS="ai,startups,fastapi"; 
python backend/app/core/landing_page/google_trends/scripts/quickstart_trends_compare.py
```

## Usage in Code
```python
from app.core.landing_page.google_trends import (
    GoogleTrendsClient, build_query, build_compare
)
from app.core.landing_page.google_trends.api._requests import Interval

client = GoogleTrendsClient()
req = build_query(terms=["dealscale"], interval=Interval.daily, last_n_days=30, region="US")

try:
    data = client.interest_over_time("dealscale", geo="US", timeframe="today 1-m")
    print(data)
except NotImplementedError:
    # Expected until a backend is wired
    print("Built request:", req.model_dump())
```

## FastAPI Integration (Dependency Injection)
```python
from fastapi import FastAPI, Depends
from app.core.landing_page.google_trends import get_trends_client_dependency, GoogleTrendsClient

app = FastAPI()

@app.get("/trends/ping")
def ping(gt: GoogleTrendsClient = Depends(get_trends_client_dependency)):
    # Use `gt` once an implementation exists
    return {"ok": True}
```

## Wiring a Backend (when ready)
Two options:

- Community: [`pytrends`](https://github.com/GeneralMills/pytrends)
  - Implement `GoogleTrendsClient.interest_over_time`, `related_queries`, and `trending_searches_daily` using pytrends calls.
  - Add `pytrends` to `pyproject.toml` (or your environment) and unit tests that mock network calls.
- Official API: If/when Google releases an official Trends API, replace method bodies with real HTTP requests using the models in `api/_requests.py` and `api/_responses.py`.

## Testing
- Placeholder unit tests exist under `_tests/` and validate builders and DI wiring.
- Add tests for retry, rate-limit, and auth handling once the transport is implemented.

## Security Notes
- Do not commit secrets (`TRENDS_API_KEY`, tokens).
- Prefer environment variables and secret stores.

## License
Apache 2.0 (match repository’s overall licensing policy if different).
