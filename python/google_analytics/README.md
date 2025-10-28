# DealScale GA4 SDK

Read-only client for Google Analytics Data API (GA4). Mirrors the structure used in the GSC SDK.

## Features
- Installed-app OAuth and Service Account auth
- Small HTTP client with retries/backoff
- Typed Pydantic request/response models
- Quickstart scripts
- Mocked tests (no real creds required)

## Install
Use your existing workspace's Python. This folder includes a `pyproject.toml` with deps similar to the GSC SDK.

## Env
Copy `.env.example` and set values, or export directly:
- `GA_CLIENT_SECRETS` (optional for installed-app)
- `GA_TOKEN_PATH` (installed-app token cache)
- `GA_SA_KEY` (service account JSON key path)
- `GA_SCOPE=https://www.googleapis.com/auth/analytics.readonly`
- `GA_PROPERTY_ID=<your numeric property id>`
- `GA_HTTP_TIMEOUT=30`
- `GA_USER_AGENT=CyberOni-GA4-SDK/1.0`

## Quickstarts
- `scripts/quickstart_ga4.py` (installed-app)
- `scripts/quickstart_ga4_service_account.py` (service account)
 - `scripts/realtime_cli.py` (dynamic realtime fetch)

## Troubleshooting
- 401/403: SA/user missing property access; wrong scope; wrong project
- Invalid argument: wrong dimension/metric; call `get_metadata()` to discover
- Quota: handle 429/5xx with backoff

## Security
- Do not commit secrets.
- Add the service account email as a user on the GA4 property in Admin.

---
 
## Setup: Google Cloud + GA4 (Service Account recommended)
 
1) Create/select a Google Cloud Project
- Visit Google Cloud Console → IAM & Admin → Create Project
- Note the Project ID
 
2) Enable the Analytics Data API
- APIs & Services → Library → search "Analytics Data API" → Enable
 
3) Create a Service Account
- IAM & Admin → Service Accounts → Create Service Account
- Assign minimal role (Viewer is sufficient for API usage); roles can be tightened later
- Create key → JSON → Download (store securely)
 
4) Grant GA4 property access to the Service Account
- In Google Analytics → Admin (gear) → Property Access Management
- Add user → use the service account email (…@<project>.iam.gserviceaccount.com)
- Give Viewer/Analyst or higher per your needs
 
5) Set environment variables (see `.env.example`)
- `GA_SA_KEY=/absolute/path/to/service-account.json`
- `GA_SCOPE=https://www.googleapis.com/auth/analytics.readonly`
- `GA_PROPERTY_ID=<your GA4 numeric property id>`
 
6) Verify access quickly
- Use the service account quickstart below to call `runReport`
 
 
## Installed-App OAuth (if you need user-consented data)
Use when each end-user must consent. Steps:
- APIs & Services → OAuth consent screen → Configure (External or Internal)
- Create Credentials → OAuth client ID → Application type: Desktop (or Web)
- Download client secrets JSON; set `GA_CLIENT_SECRETS` and `GA_TOKEN_PATH`
- First run prompts the user; refresh token is stored and used automatically later
 
 
## Quickstart: Service Account
File: `backend/app/core/landing_page/google_analytics/scripts/quickstart_ga4_service_account.py`
 
```python
from app.core.landing_page.google_analytics import GA4Client, build_run_report_request
import os
 
property_id = os.environ["GA_PROPERTY_ID"]
client = GA4Client.from_service_account(key_path=os.environ["GA_SA_KEY"])  # uses readonly scope by default
 
req = build_run_report_request(
    dimensions=["date"],
    metrics=["activeUsers"],
    last_n_days=7,
    limit=10,
)
resp = client.run_report(property_id, req)
print(resp.rowCount, [r.dimensionValues for r in resp.rows])
```
 
 
## Quickstart: Installed-App OAuth
File: `backend/app/core/landing_page/google_analytics/scripts/quickstart_ga4.py`
 
```python
import os
from app.core.landing_page.google_analytics import GA4Client, build_run_report_request
 
client = GA4Client.from_installed_app(
    client_secrets_file=os.environ["GA_CLIENT_SECRETS"],
    token_cache_path=os.environ.get("GA_TOKEN_PATH", ".ga4_token.json"),
)
 
req = build_run_report_request(
    dimensions=["country"], metrics=["activeUsers"], last_n_days=3
)
resp = client.run_report(os.environ["GA_PROPERTY_ID"], req)
print(resp.rowCount)
```
 
 
## Using the SDK in FastAPI
 
You can wire the client into routes using simple utilities.
 
- Dependency that builds a `GA4Client` from a Bearer token header: `get_ga4_client_dependency` in `google_analytics/__init__.py` (from `api/deps.py`).
- Convenience request builders: `build_run_report_request`, `build_realtime_request`.
 
Example route usage:
 
```python
# file: app/api/routes/analytics_example.py
from fastapi import APIRouter, Depends
from app.core.landing_page.google_analytics import (
    get_ga4_client_dependency, build_run_report_request, GA4Client,
)
 
router = APIRouter(prefix="/analytics", tags=["analytics"])
 
@router.get("/active-users")
def active_users_last_7_days(ga: GA4Client = Depends(get_ga4_client_dependency)):
    req = build_run_report_request(dimensions=["date"], metrics=["activeUsers"], last_n_days=7)
    resp = ga.run_report("<PROPERTY_ID>", req)
    return resp.model_dump()
```

## Realtime: Dynamic Fetch

Use the realtime CLI to query any dimensions/metrics supported by GA4 realtime:

```bash
python -m google_analytics.scripts.realtime_cli \
  --property "$GA_PROPERTY_ID" \
  --metrics activeUsers,eventCount \
  --dimensions eventName,streamId \
  --limit 10
```

Auth switches to installed-app with `--installed-app`. By default, the CLI uses service account credentials from `GA_SA_KEY`.
 
If you prefer service accounts in the API layer (no user token):
 
```python
from app.core.landing_page.google_analytics import GA4Client, build_run_report_request
import os
 
ga = GA4Client.from_service_account(key_path=os.environ["GA_SA_KEY"])  # initialize once at startup
 
def handler():
    req = build_run_report_request(dimensions=["pagePath"], metrics=["screenPageViews"], last_n_days=1)
    return ga.run_report(os.environ["GA_PROPERTY_ID"], req).model_dump()
```
 
 
## Troubleshooting Cheatsheet
- 401/403
  - Service account lacks GA4 property access, or user token is invalid
  - Wrong Cloud project or API not enabled
- Invalid argument
  - Dimension/metric name incorrect; call `get_metadata()` to inspect available fields
- Quotas / 429 / 5xx
  - Retries with backoff are built-in; tune timeouts if needed
