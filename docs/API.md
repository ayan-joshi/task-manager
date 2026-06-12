# API Reference

Base URL: `${NEXT_PUBLIC_API_URL}` → `http://localhost:8080/api` in local development.

## Conventions

### Response envelope

Every response uses a consistent envelope:

```jsonc
// success
{ "success": true, "data": { /* ... */ }, "meta": { /* pagination, when applicable */ } }

// error
{ "success": false, "error": { "code": "TASK_NOT_FOUND", "message": "task not found" } }
```

### Authentication

Obtain a token via `/auth/signup` or `/auth/login`, then send it on every protected request:

```
Authorization: Bearer <token>
```

### Status codes

| Code | Meaning |
| ---- | ------- |
| 200 | OK |
| 201 | Created |
| 401 | Missing/invalid token or bad credentials |
| 403 | Authenticated but not permitted (RBAC) |
| 404 | Resource not found (also returned for cross-user access) |
| 409 | Conflict (email already registered) |
| 413 | Uploaded file too large |
| 415 | Unsupported file type |
| 422 | Validation error (includes `details[]`) |
| 429 | Rate limited |

### Enums

- **status:** `TODO`, `IN_PROGRESS`, `COMPLETED`
- **priority:** `LOW`, `MEDIUM`, `HIGH`
- **role:** `USER`, `ADMIN`

---

## Auth

### POST `/auth/signup`

```json
{ "name": "Ada Lovelace", "email": "ada@example.com", "password": "supersecret" }
```

`201` → `{ "success": true, "data": { "token": "...", "user": { ... } } }`

### POST `/auth/login`

```json
{ "email": "ada@example.com", "password": "supersecret" }
```

`200` → `{ "token": "...", "user": { ... } }`

### GET `/auth/me` 🔒

`200` → the authenticated user.

---

## Tasks 🔒

### GET `/tasks`

Query parameters (all optional, all combinable):

| Param | Type | Notes |
| ----- | ---- | ----- |
| `search` | string | case-insensitive title match |
| `status` | enum | filter by status |
| `priority` | enum | filter by priority |
| `sort` | enum | `created_desc` (default), `created_asc`, `due_date_asc`, `due_date_desc`, `priority_asc`, `priority_desc`, `title_asc`, `title_desc` |
| `scope` | enum | `me` (default) returns only your tasks. `all` returns every user's tasks **and** includes each task's `owner` — honored for `ADMIN` only; silently falls back to `me` for regular users. |
| `page` | int | default `1` |
| `limit` | int | default `10`, max `100` |

`200`:

```jsonc
{
  "success": true,
  "data": [ { /* task */ } ],
  "meta": { "total": 42, "page": 1, "limit": 10, "totalPages": 5 }
}
```

### POST `/tasks`

```json
{
  "title": "Write report",
  "description": "Q3 numbers",
  "status": "TODO",
  "priority": "HIGH",
  "dueDate": "2026-07-01T00:00:00Z"
}
```

`201` → the created task. Only `title` is required.

### GET `/tasks/:id`

`200` → the task (with `attachments`). `404` if not found or not owned.

### PUT `/tasks/:id`

Partial update — every field optional:

```json
{ "status": "COMPLETED", "clearDueDate": true }
```

`clearDueDate: true` removes the due date. `200` → the updated task.

### DELETE `/tasks/:id`

`200` → `{ "deleted": true, "id": "..." }`

### GET `/tasks/:id/activity`

`200` → array of activity log entries for the task (paginated via `page`/`limit`).

---

## Attachments 🔒

### POST `/tasks/:id/attachments`

`multipart/form-data` with a `file` field. Validates against `MAX_UPLOAD_MB` and `ALLOWED_MIME_TYPES`.

`201` → attachment metadata.

### GET `/attachments/:attachmentId`

Streams the file with `Content-Disposition: attachment`. (The frontend fetches this with the bearer header and triggers a browser download.)

### DELETE `/attachments/:attachmentId`

`200` → `{ "deleted": true, "id": "..." }`

---

## Real-time — GET `/events`

Server-Sent Events stream. Because `EventSource` cannot send headers, the token is passed as a query parameter:

```
GET /api/events?token=<jwt>
```

Emits named events scoped to the authenticated user:

```
event: task.created
data: { "type": "task.created", "taskId": "...", "task": { ... } }
```

Event types: `task.created`, `task.updated`, `task.deleted`.

---

## Admin 🔒 (role `ADMIN`)

### GET `/admin/users`

`200` → paginated list of all users.

### GET `/admin/activity-logs`

`200` → paginated global activity log, each entry preloaded with its `user`.

---

## Health — GET `/health`

`200` → `{ "success": true, "data": { "status": "ok" } }` (no auth required).
