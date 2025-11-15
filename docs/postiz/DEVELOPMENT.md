# Postiz Development Guide

This page explains how to bootstrap a local environment, understand the
architecture, and use the key tooling inside the Postiz monorepo.

---

## 1. Prerequisites

- Node.js 18+ (PNPM or npm, depending on preference)
- Docker / Docker Compose (Postgres + Redis, optional if you run managed
  services)
- `nx` CLI (installed via `npx` or globally)
- A `.env` file at the root (shared across all apps)

Clone the repo and install dependencies:

```bash
git clone https://github.com/postizapp/postiz.git
cd postiz
npm install
```

Copy the sample environment file and tweak values for local development:

```bash
cp .env.example .env
# Edit database / redis / OAuth secrets, etc.
```

Run database migrations and generate the Prisma client:

```bash
npm run prisma-generate
npm run prisma-db-push
```

---

## 2. Architecture Overview

- **Frontend**: Next.js + Tailwind, lives under `apps/frontend`.
- **Backend**: NestJS app (`apps/backend`) using controllers/services/repos/DTOs.
  Prisma provides database access, Redis handles scheduling and background jobs.
- **Cron**: NestJS app (`apps/cron`) sharing code with the backend to run
  recurring tasks.
- **Worker**: NestJS app (`apps/worker`) responsible for async jobs and post
  execution.
- **NX Monorepo**: The entire repo is orchestrated via NX; scripts are defined in
  the root `package.json`.

All apps share a single `.env` file for simplicity. This is intentionally
different from the usual NX-per-app configuration.

---

## 3. Common Scripts

Run these from the repository root unless noted otherwise:

| Command | Description |
| --- | --- |
| `npm run dev` | Starts the development server (frontend + backend via NX). |
| `npm run prisma-generate` | Generates the Prisma client. |
| `npm run prisma-db-push` | Pushes the current Prisma schema to the configured DB (useful for rapid prototyping). |
| `nx graph` | Visualize project dependencies. |

Individual apps can also be run with NX targets, e.g.:

```bash
npx nx serve backend
npx nx serve frontend
npx nx serve cron
npx nx serve worker
```

---

## 4. Frontend Setup

Located in `apps/frontend` (Next.js + Tailwind). The frontend consumes the same
`.env` file and expects `NEXT_PUBLIC_BACKEND_URL` to point to the backend server.

During local development:

1. Ensure `FRONTEND_URL` and `NEXT_PUBLIC_BACKEND_URL` in `.env` point to your
   local hostnames (e.g., `http://localhost:4200`, `http://localhost:3000`).
2. Run `npm run dev` to start the NX orchestrated dev server, or target the app
   directly via `npx nx serve frontend`.

---

## 5. Backend / Cron / Worker

All backend services are NestJS apps under `apps/backend`, `apps/cron`, and
`apps/worker`. They share DTOs, providers, and infrastructure code via
`libraries/nestjs-libraries`.

Key services:

- **Backend**: API, auth, scheduling endpoints, Prisma access.
- **Cron**: Scheduled jobs (daily syncs, periodic cleanups).
- **Worker**: Handles queued jobs (posting content, refreshing tokens, etc.).

Prerequisites:

- Postgres database (default connection string in `.env` is Postgres, but Prisma
  supports MySQL/MariaDB if you change the provider and connection string).
- Redis instance for queues and rate limiting.

Start the backend with `npx nx serve backend`. Cron/worker can be started in
separate terminals using the corresponding NX commands.

---

## 6. Contributors Guide

The repository includes a contributor guide detailing pull request conventions,
branch naming, and code style. In short:

- Fork → feature branch → PR against `main`.
- Include tests where possible.
- Follow the existing TypeScript/Prisma conventions for DTOs and migrations.

Refer to `CONTRIBUTING.md` in the root for the full guidelines.

---

## 7. Tips

- **Shared `.env`**: Remember that all apps read from the same `.env`. Update
  it whenever you add provider credentials, OAuth settings, or storage configs.
- **Prisma**: Use `npx prisma studio` to inspect the local database.
- **Redis**: Flush queues between runs if you need a clean slate (`redis-cli
  FLUSHALL`).
- **Docker**: If you don’t want to manage Postgres/Redis manually, use the
  compose stack under `frontend/postiz` or your own docker compose file to spin
  them up quickly.

Happy hacking!


