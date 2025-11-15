# Postiz Frontend Environment

This folder hosts a self-contained Docker Compose stack for running the
open-source Postiz app locally. It mirrors the reference deployment described in
the upstream docs and is tuned for VM-class hardware (2 vCPU, 2 GB RAM).

## Structure

| File | Purpose |
| --- | --- |
| `docker-compose.yml` | Services for Postiz, Postgres, and Redis on a shared network. |
| `postiz.env.example` | Example configuration values you can copy to `postiz.env`. |

## Quick Start

1. Copy the example env file and update values as needed:
   ```bash
   cd frontend/postiz
   cp postiz.env.example postiz.env
   # edit postiz.env with your domain, secrets, OAuth keys, etc.
   ```

2. Start the stack:
   ```bash
   docker compose --env-file postiz.env up -d
   ```

3. (Optional) Stop / remove containers:
   ```bash
   docker compose --env-file postiz.env down
   ```

## Environment Variables

All configuration is provided via the `postiz.env` file. The most important
settings are:

| Variable | Description |
| --- | --- |
| `MAIN_URL`, `FRONTEND_URL` | Public base URL where users access Postiz. |
| `NEXT_PUBLIC_BACKEND_URL` | Public API endpoint (usually `${MAIN_URL}/api`). |
| `JWT_SECRET` | Random string unique per deployment. |
| `DATABASE_URL` | Connection string to the Postgres service. |
| `REDIS_URL` | Connection string to the Redis service. |
| `NOT_SECURED` | Set to `true` only for non-HTTPS development environments. |
| `POSTIZ_APPS` | Control which apps run inside the container (`frontend`, `backend`, `worker`, `cron`). Leave empty to run all. |

For full API / OAuth / provider keys, consult the upstream configuration
reference. Any value defined in `postiz.env` is automatically passed into the
Postiz container via `env_file`.

## Networking Notes

- Only port `5000` is published by default. Route HTTPS traffic from a reverse
  proxy (Caddy, Nginx, Traefik) to this port.
- Internal services (frontend `4200`, backend `3000`, Postgres `5432`, Redis
  `6379`) stay on the private `postiz-network`.
- Set `NOT_SECURED=true` only if you cannot provide HTTPS and you accept the
  security trade-off (cookies marked non-secure).

## Updating Configuration

Whenever you change environment variables, run:
```bash
docker compose --env-file postiz.env down
docker compose --env-file postiz.env up -d
```

This recreates the containers with the new settings.

## Configuration Patterns

You will encounter provider-specific variables such as:

```
INSTAGRAM_CLIENT_ID=1234567890
LINKEDIN_CLIENT_SECRET=abcdef
```

You can inject them in three different ways:

1. **Inline in `docker-compose.yml`:**
   ```yaml
   services:
     postiz:
       environment:
         INSTAGRAM_CLIENT_ID: "1234567890"
         INSTAGRAM_CLIENT_SECRET: "foo"
   ```

2. **From an env file (`postiz.env` in this repo):**
   ```yaml
   services:
     postiz:
       env_file:
         - postiz.env
   ```
   ```
   # postiz.env
   INSTAGRAM_CLIENT_ID=1234567890
   INSTAGRAM_CLIENT_SECRET=foo
   ```

3. **Hybrid:** keep defaults inline and override the rest via env file.
   ```yaml
   services:
     postiz:
       environment:
         INSTAGRAM_CLIENT_ID: "1234567890"
       env_file:
         - postiz.env   # contains INSTAGRAM_CLIENT_SECRET, etc.
   ```

**Important:** when you rely on `env_file`, Docker Compose expects all values to
come from that file. In other words, every variable defined in the env file
overrides the inline definitions. For shared services such as Postgres/Redis, it
is fine to keep their credentials in the compose file, but app-specific secrets
should live inside `postiz.env`.

## Email Providers

Postiz supports two outbound email options:

### Resend
```env
EMAIL_PROVIDER=resend
EMAIL_FROM_NAME="Postiz Emailer"
EMAIL_FROM_ADDRESS="postiz@example.com"
RESEND_API_KEY=pk_xxx
```
1. Create a Resend account and verify your domain.
2. Copy the API key and place it in `postiz.env`.
3. When Resend is enabled, new users must confirm their email address before logging in.

### NodeMailer (SMTP)
```env
EMAIL_PROVIDER=nodemailer
EMAIL_FROM_NAME="Postiz Emailer"
EMAIL_FROM_ADDRESS="postiz@example.com"
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=465
EMAIL_SECURE=true
EMAIL_USER=username
EMAIL_PASS=password
```
This mode works with any SMTP server (Gmail, SendGrid, SES, etc.). Set the host,
port, TLS flag, and credentials accordingly. Activation emails are also enforced
when SMTP is configured.

## Cloudflare R2 Storage (Optional)

If you prefer not to keep uploads locally:

1. Log in to the Cloudflare dashboard → R2 Object Storage → create a bucket
   (Automatic / Standard).
2. Under **API Tokens**, create a token with “Object Read & Write” access scoped
   to your bucket. Copy the **Account ID**, **Access Key ID**, and **Secret
   Access Key**.
3. Connect a custom domain (or use the default `*.r2.cloudflarestorage.com`
   endpoint) and configure a permissive CORS policy, e.g.:
   ```json
   [
     {
       "AllowedOrigins": [
         "http://localhost:4200",
         "https://yourDomain.com"
       ],
       "AllowedMethods": ["GET", "POST", "HEAD", "PUT", "DELETE"],
       "AllowedHeaders": [
         "Authorization",
         "x-amz-date",
         "x-amz-content-sha256",
         "content-type"
       ],
       "ExposeHeaders": ["ETag", "Location"],
       "MaxAgeSeconds": 3600
     }
   ]
   ```
4. Populate `postiz.env`:
   ```env
   STORAGE_PROVIDER=cloudflare
   CLOUDFLARE_ACCOUNT_ID=acc_xxx
   CLOUDFLARE_ACCESS_KEY=AKIA...
   CLOUDFLARE_SECRET_ACCESS_KEY=...
   CLOUDFLARE_BUCKETNAME=postiz-media
   CLOUDFLARE_BUCKET_URL=https://uploads.example.com
   CLOUDFLARE_REGION=wnam
   ```
   Set `NEXT_PUBLIC_UPLOAD_DIRECTORY` to match your bucket URL if you want
   Pre-signed URLs accessible directly via the CDN.

## OAuth / OIDC Login

To enable login via an external identity provider (Authentik, Keycloak, Dex, etc.):

1. Create an application on the IdP with:
   - `redirect_uri`: `https://postiz.yourserver.com/settings`
   - Client ID / secret, plus authorization, token, and userinfo endpoints.
2. Populate `postiz.env`:
   ```env
   POSTIZ_GENERIC_OAUTH=true
   NEXT_PUBLIC_POSTIZ_OAUTH_DISPLAY_NAME=Authentik
   NEXT_PUBLIC_POSTIZ_OAUTH_LOGO_URL=https://raw.githubusercontent.com/walkxcode/dashboard-icons/master/png/authentik.png
   POSTIZ_OAUTH_URL=https://authentik.example.com
   POSTIZ_OAUTH_AUTH_URL=https://authentik.example.com/application/o/authorize
   POSTIZ_OAUTH_TOKEN_URL=https://authentik.example.com/application/o/token
   POSTIZ_OAUTH_USERINFO_URL=https://authentik.example.com/application/o/userinfo
   POSTIZ_OAUTH_CLIENT_ID=randomclientid
   POSTIZ_OAUTH_CLIENT_SECRET=randomclientsecret
   ```
3. Restart the stack. The login screen will show a “Sign in with Authentik”
   button using the display name/logo you provided.

## Required vs Optional Settings

Postiz reads **all** configuration from environment variables. Any change means
you must restart the containers for it to take effect (`docker compose down &&
docker compose up -d`). The most important keys are already pre-populated in
`postiz.env.example`:

| Variable | Description |
| --- | --- |
| `DATABASE_URL` | Prisma connection string (PostgreSQL by default). You can point to MySQL/MariaDB as long as Prisma supports the driver. |
| `REDIS_URL` | Connection string to Redis (used for queues + rate limiting). |
| `JWT_SECRET` | Random string used to sign auth tokens; generate a unique value per deployment. |
| `FRONTEND_URL` | Public URL where users access the UI (e.g., `https://postiz.example.com`). |
| `NEXT_PUBLIC_BACKEND_URL` | Public API endpoint (typically `${FRONTEND_URL}/api`). |
| `BACKEND_INTERNAL_URL` | Internal API endpoint when running all services locally (`http://localhost:3000`). |

Optional flags (see `postiz.env.example` for defaults):

- `DISABLE_REGISTRATION=true` – only the first user can sign up; afterwards the
  registration page is hidden (also disables OAuth).
- `DISABLE_IMAGE_COMPRESSION=true` – skip server-side compression if you need
  original asset quality.
- Provider/OAuth/storage/email variables covered in sections above.

Whenever you add new social media credentials (e.g., LinkedIn, X, Threads),
append them to `postiz.env`, restart the stack, and re-run the provider
connection flow inside Postiz.


