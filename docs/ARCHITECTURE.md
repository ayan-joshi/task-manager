# Architecture

## Overview

TaskFlow is a decoupled two-tier application. The Go API owns all business logic, authorization, and persistence; the Next.js client is a pure consumer of the HTTP + SSE API. They share no code and deploy independently (Railway + Vercel), communicating only over the network.

```
┌─────────────────────┐         HTTP / SSE          ┌──────────────────────┐
│   Next.js frontend  │  ───────────────────────▶   │      Go API (Gin)     │
│  (App Router, RSC)  │  ◀───────  JSON / events ─  │  handler→service→repo │
└─────────────────────┘                              └───────────┬──────────┘
                                                                  │ GORM
                                                       ┌──────────▼──────────┐
                                                       │     PostgreSQL       │
                                                       └──────────────────────┘
```

## Backend — clean layering

Requests flow through clearly separated layers, each depending only on the one beneath it via interfaces:

```
HTTP request
   │
   ▼
middleware (request-id → log → recover → security headers → CORS → rate limit → [auth → rbac])
   │
   ▼
handler        // parse & validate input, map results to the response envelope
   │
   ▼
service        // business rules, authorization (ownership), orchestration
   │           //   → publishes SSE events, records activity log
   ▼
repository     // data access only (GORM); no business logic
   │
   ▼
PostgreSQL
```

### Why this shape

- **Testability.** Services depend on repository *interfaces*, and handlers are driven through the real router in integration tests. The whole stack is exercised against in-memory SQLite, so tests are fast and need no external services.
- **Separation of concerns.** Validation lives at the edge (DTO binding tags), authorization lives in services, and SQL lives in repositories. A handler never builds a query; a repository never decides who may see a row.
- **Explicit dependency injection.** `server.NewContainer` wires every repository, service, and handler once. The dependency graph is visible in a single file rather than hidden behind a framework or globals.

### Key decisions

| Concern | Decision | Rationale |
| ------- | -------- | --------- |
| Identifiers | UUIDs generated in a GORM `BeforeCreate` hook | DB-agnostic (works on Postgres & SQLite); IDs are non-enumerable |
| Authorization | Ownership checked in the service; cross-user access returns `404` | Avoids disclosing existence of other users' data |
| Filtering/sort | Done entirely in SQL with a hard-coded `ORDER BY` map | No in-memory filtering; the sort clause never contains user input (injection-safe) |
| Search | `LOWER(title) LIKE ? ESCAPE '\'` with wildcard escaping | Portable across Postgres/SQLite; literal `%`/`_` don't match everything |
| Real-time | In-process SSE broker, per-user channels | Simple, no external broker; same ownership isolation as REST |
| Files | `Storage` interface with a local-disk implementation | Swap to S3/GCS without touching services or handlers |
| Activity log | Best-effort write (logged, never fatal) | An audit failure must not break the user's operation |
| Errors | Typed `AppError` (status + code + message) | Consistent envelope; services stay HTTP-agnostic |

### Security

- bcrypt password hashing; passwords never serialized (`json:"-"`).
- HS256 JWTs with enforced signing method (prevents `alg` confusion).
- Per-IP token-bucket rate limiting with idle eviction.
- Security headers (CSP, `X-Frame-Options`, `nosniff`, etc.), configurable CORS allow-list.
- Parameterized queries throughout (GORM); directory-traversal-safe file storage.
- Panic recovery middleware converts crashes into clean `500`s.
- Graceful shutdown drains in-flight requests on `SIGINT`/`SIGTERM`.

## Frontend — App Router + server-state

- **Routing:** App Router pages under `src/app`. Authenticated screens are wrapped by `AuthGuard` (waits for persisted auth to rehydrate, then redirects) and `AppShell` (sidebar, header, user menu, theme toggle).
- **Server state:** TanStack Query owns all API data. Mutations implement **optimistic updates with rollback** via `onMutate`/`onError`/`onSettled`, and invalidate on settle for authoritative reconciliation.
- **Client state:** Zustand persists `{ token, user }` to `localStorage` so refresh keeps the session. An Axios interceptor attaches the bearer token and clears auth on `401`.
- **Real-time:** `useTaskEvents` opens an `EventSource` to `/events` and invalidates task queries on push, so changes from other sessions appear without polling.
- **Forms & validation:** React Hook Form + Zod resolvers; the same enums and constraints mirror the backend's validation.
- **UI:** shadcn/ui (Radix primitives) + Tailwind, dark mode via `next-themes`, toasts via sonner, accessible labels and ARIA throughout.

## Data model

```
User 1 ──── * Task 1 ──── * Attachment
  │                │
  └──── * ActivityLog * ─┘   (ActivityLog references user_id and optional task_id)
```

- `users(id, name, email[unique], password_hash, role, timestamps)`
- `tasks(id, user_id[fk], title, description, status, priority, due_date, timestamps)`
- `attachments(id, task_id[fk], user_id, file_name, stored_name, content_type, size_bytes, created_at)`
- `activity_logs(id, user_id, task_id?, action, metadata(json), created_at)`

Cascade deletes (`OnDelete:CASCADE`) ensure a deleted task removes its attachments, and a deleted user removes their tasks.

## CI/CD

GitHub Actions runs two parallel jobs on push/PR:

- **backend:** `go vet`, `go test -race`, `go build`.
- **frontend:** `npm ci`, `npm run lint`, `npm run test`, `npm run build`.

Containers build from per-app multi-stage Dockerfiles; `docker-compose.yml` brings up Postgres + backend + frontend for one-command local parity with production.
