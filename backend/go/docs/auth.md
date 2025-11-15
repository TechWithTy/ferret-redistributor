# Authentication

Two token strategies are supported:

- JWT (HS256)
  - Enabled when `JWT_SECRET` is set.
  - API issues JWT with `sub` (user id), `iat`, `exp`, optional `iss`.
  - Middleware validates signature and expiry.

- Opaque sessions (default)
  - API creates a random token and stores its SHA‑256 hash in `auth_sessions` with expiry.
  - Middleware looks up the token hash and checks expiry/revocation.

Flows
- Signup (email/password):
  - Creates `users` row and `auth_identities` with bcrypt hash (build with `-tags=secure`).
  - Optionally send email verification using `auth_email_verifications`.
- Login (email/password):
  - Verifies credentials; issues JWT or session token.
- Forgot/Reset:
  - Issues one‑time token in `auth_password_resets`; on reset, rotates password hash.
- Phone verification:
  - Uses `auth_phone_codes` with hashed OTPs; updates `verified_at` on phone identity.
- OAuth (LinkedIn/Meta):
  - Stores provider id in `identifier` and tokens in `oauth_data` for `auth_identities`.

