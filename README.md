# TaskFlow — Full-Stack Task Management

A production-grade task management application: a **Go (Gin + GORM)** REST API backed by **PostgreSQL**, and a **Next.js 15 (App Router) + TypeScript** frontend. Includes JWT auth with per-user task isolation, role-based admin access, real-time updates over SSE, optimistic UI, file attachments, an activity audit log, dark mode, Docker, and CI.

> Monorepo: `apps/backend` (Go) and `apps/frontend` (Next.js), orchestrated with Docker Compose.

---

## Table of contents

- [Features](#features)
- [Tech stack](#tech-stack)
- [Project structure](#project-structure)
- [Quick start (Docker)](#quick-start-docker)
- [Local development](#local-development)
- [Environment variables](#environment-variables)
- [Testing](#testing)
- [API overview](#api-overview)
- [Deployment](#deployment-railway--vercel)
- [Assumptions & trade-offs](#assumptions--trade-offs)

---

## Features

**Core**

- JWT authentication (signup / login), bcrypt password hashing, persisted session (survives refresh).
- Task CRUD with server-side validation and consistent error responses.
- Filtering by status & priority, title search, and multi-field sort — **all combinable** — with pagination metadata.
- Per-user isolation: users only ever see and mutate their own tasks (enforced in the service layer).
- Responsive UI (mobile → desktop) with table and card views, loading skeletons, empty/error states, and toasts.

**Bonus (all implemented)**

- **Role-based access** — `ADMIN` role can list all users and view the global activity log.
- **Real-time updates** — Server-Sent Events push task changes to connected clients instantly.
- **Optimistic UI** — create / update / delete apply immediately and roll back on failure (TanStack Query).
- **File attachments** — upload, download, and delete files on a task with type/size validation.
- **Activity log** — every task change is recorded and shown on a per-task timeline (and admin-wide).
- **Dockerized** — `docker compose up --build` launches Postgres + backend + frontend.
- **CI** — GitHub Actions runs lint, tests, and builds for both apps.
- **Dark mode** — theme toggle with persisted preference.

---

## Tech stack

| Layer    | Technology |
| -------- | ---------- |
| Backend  | Go 1.22, Gin, GORM, PostgreSQL, golang-jwt, bcrypt, validator |
| Frontend | Next.js 15 (App Router), TypeScript, Tailwind CSS, shadcn/ui (Radix), TanStack Query, React Hook Form, Zod, Zustand, Axios |
| Infra    | Docker, Docker Compose, GitHub Actions |

Architecture: clean layering (handler → service → repository), dependency injection via an explicit container, repository pattern, and a dedicated middleware chain. See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

---

## Project structure

```
task-manager/
├── apps/
│   ├── backend/                 # Go API
│   │   ├── cmd/api/main.go       # entrypoint (graceful shutdown)
│   │   ├── internal/
│   │   │   ├── config/           # env loading + validation
│   │   │   ├── domain/           # GORM models & enums
│   │   │   ├── database/         # connection + migrations
│   │   │   ├── dto/              # request/response types
│   │   │   ├── repository/       # data access (repository pattern)
│   │   │   ├── service/          # business logic + authorization
│   │   │   ├── handler/          # HTTP handlers
│   │   │   ├── middleware/       # auth, RBAC, CORS, rate limit, security
│   │   │   ├── sse/              # real-time event broker
│   │   │   ├── storage/          # pluggable file storage
│   │   │   └── server/           # DI container, routes, seeding
│   │   ├── tests/                # API integration tests (SQLite)
│   │   └── Dockerfile
│   └── frontend/                # Next.js app
│       ├── src/app/             # routes: /login /signup /dashboard /tasks/[id] /admin
│       ├── src/components/      # UI primitives + feature components
│       ├── src/hooks/           # TanStack Query hooks, SSE, auth
│       ├── src/lib/             # api client, types, validations, store
│       ├── src/__tests__/       # Vitest component tests
│       └── Dockerfile
├── docs/                        # API.md, ARCHITECTURE.md
├── .github/workflows/ci.yml
├── docker-compose.yml
└── .env.example
```

---

## Quick start (Docker)

Prerequisites: Docker + Docker Compose.

```bash
git clone <repo-url>
cd task-manager
cp .env.example .env        # optional: adjust JWT_SECRET / seed admin

docker compose up --build
```

Then open:

- **Frontend:** http://localhost:3000
- **Backend health:** http://localhost:8080/api/health

A demo admin account is seeded on first boot (defaults `admin@example.com` / `admin12345`, configurable in `.env`). Create your own user via the signup page.

To stop and remove volumes:

```bash
docker compose down -v
```

---

## Local development

You need Go 1.22+, Node 22+, and a PostgreSQL instance.

### 1. Database

Either use the compose Postgres only:

```bash
docker compose up -d postgres
```

or point `DATABASE_URL` at any Postgres you control.

### 2. Backend

```bash
cd apps/backend
cp .env.example .env          # set DATABASE_URL + JWT_SECRET
go mod download
go run ./cmd/api              # serves on :8080, auto-migrates on boot
```

### 3. Frontend

```bash
cd apps/frontend
cp .env.example .env.local    # NEXT_PUBLIC_API_URL=http://localhost:8080/api
npm install
npm run dev                   # serves on :3000
```

---

## Environment variables

### Backend (`apps/backend/.env`)

| Variable | Required | Default | Description |
| -------- | -------- | ------- | ----------- |
| `DATABASE_URL` | ✅ | — | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | — | JWT signing secret (≥ 16 chars) |
| `JWT_EXPIRY_HOURS` | | `72` | Token lifetime |
| `PORT` | | `8080` | HTTP port |
| `ALLOWED_ORIGINS` | | `http://localhost:3000` | Comma-separated CORS origins |
| `UPLOAD_DIR` | | `./uploads` | Attachment storage directory |
| `MAX_UPLOAD_MB` | | `10` | Max attachment size |
| `ALLOWED_MIME_TYPES` | | (images, pdf, text) | Allowed attachment content types |
| `RATE_LIMIT_RPS` / `RATE_LIMIT_BURST` | | `10` / `20` | Per-IP rate limit |
| `ENVIRONMENT` | | `development` | `development` or `production` |
| `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` | | — | Optional admin seed |

### Frontend (`apps/frontend/.env.local`)

| Variable | Required | Description |
| -------- | -------- | ----------- |
| `NEXT_PUBLIC_API_URL` | ✅ | Base URL of the backend API, e.g. `http://localhost:8080/api` |

---

## Testing

**Backend** — integration tests run the full HTTP stack against an in-memory SQLite database (no Docker required):

```bash
cd apps/backend
go test ./...          # auth, task CRUD + ownership, filtering/sort/pagination
```

**Frontend** — component tests with Vitest + React Testing Library:

```bash
cd apps/frontend
npm run test           # login form, task list, create-task form
```

---

## API overview

Base path: `/api`. All task/admin routes require `Authorization: Bearer <token>`. Full reference in [`docs/API.md`](docs/API.md).

| Method | Path | Description |
| ------ | ---- | ----------- |
| POST | `/auth/signup` | Register |
| POST | `/auth/login` | Authenticate |
| GET | `/auth/me` | Current user |
| GET | `/tasks` | List (filter/search/sort/paginate) |
| POST | `/tasks` | Create |
| GET | `/tasks/:id` | Fetch one |
| PUT | `/tasks/:id` | Update |
| DELETE | `/tasks/:id` | Delete |
| GET | `/tasks/:id/activity` | Task activity log |
| POST | `/tasks/:id/attachments` | Upload attachment |
| GET | `/attachments/:id` | Download attachment |
| DELETE | `/attachments/:id` | Delete attachment |
| GET | `/events` | SSE real-time stream |
| GET | `/admin/users` | List users (admin) |
| GET | `/admin/activity-logs` | Global activity log (admin) |

---

## Deployment (Railway + Vercel)

### Backend → Railway

1. Create a Railway project and add a **PostgreSQL** plugin.
2. Add a service from this repo with root directory `apps/backend` (it has a `Dockerfile`).
3. Set variables: `DATABASE_URL` (reference the Postgres plugin), `JWT_SECRET`, `ENVIRONMENT=production`, `ALLOWED_ORIGINS=https://<your-vercel-domain>`, and optionally `SEED_ADMIN_EMAIL`/`SEED_ADMIN_PASSWORD`.
4. Deploy. Migrations run automatically on boot. Note the public URL, e.g. `https://taskflow-api.up.railway.app`.

> Attachments use local disk by default. On Railway, attach a **Volume** mounted at `/app/uploads` to persist files across deploys (the storage layer is an interface, so swapping in S3/GCS is a drop-in change).

### Frontend → Vercel

1. Import the repo into Vercel; set the **Root Directory** to `apps/frontend`.
2. Set environment variable `NEXT_PUBLIC_API_URL=https://<your-railway-domain>/api`.
3. Deploy (framework auto-detected as Next.js).
4. Add the resulting Vercel domain to the backend's `ALLOWED_ORIGINS`.

---

## Assumptions & trade-offs

- **Auto-migration over SQL migration files.** GORM `AutoMigrate` keeps the schema in sync on boot — fast for an assessment. For stricter production change control, these would be exported to versioned migrations.
- **SQLite for tests, Postgres for runtime.** Models are written database-agnostically (UUIDs generated in Go, portable `LIKE ... ESCAPE` search), so tests run anywhere without Docker while production uses Postgres.
- **In-memory rate limiter / SSE broker.** Sufficient for single-node deployment. A horizontally-scaled deployment would back these with Redis (rate limit) and a pub/sub fan-out (SSE).
- **Local-disk attachments behind a `Storage` interface.** Swapping to object storage requires implementing one interface, no call-site changes.
- **`PUT` for updates.** The brief mentioned `PATCH`; the API performs partial updates (all fields optional) and exposes them under `PUT` per the detailed spec. Behaviour is partial-update semantics either way.
- **404 instead of 403 on cross-user access.** Accessing another user's task returns `404` so the existence of others' resources is never disclosed.
