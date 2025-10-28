# API

- Base URL defaults to `http://localhost:8080`.
- Static spec: GET `/openapi.json`.

Auth endpoints
- POST `/v1/auth/signup` — body: email, password, display_name, org_id (optional); returns user_id.
- POST `/v1/auth/login` — body: email, password; returns bearer token.
- POST `/v1/auth/forgot` — body: email; always returns ok.
- POST `/v1/auth/reset` — body: token, new_password.

Auth
- If `JWT_SECRET` is set, the API issues and accepts HS256 JWTs.
- Otherwise, it issues opaque session tokens stored in `auth_sessions`.
- Send bearer token via `Authorization: Bearer <token>`.

User/profile
- GET `/v1/users/:id` — demo payload.
- GET `/v1/profile` — returns user profile (or null if missing).
- PUT `/v1/profile` — upserts user profile (display, timezone, locale, prefs).

ICP (org personalization)
- GET `/v1/icp` — fetches first ICP profile for user’s org.
- PUT `/v1/icp` — upserts ICP profile by `(org_id, name)`; defaults to `Default`.

