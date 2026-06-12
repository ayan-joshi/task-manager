# Docker

Container build files live alongside each app:

- `apps/backend/Dockerfile` — multi-stage build producing a static Go binary on a minimal Alpine runtime, running as a non-root user with a healthcheck.
- `apps/frontend/Dockerfile` — multi-stage build producing a Next.js standalone server, running as a non-root user.

The root [`docker-compose.yml`](../docker-compose.yml) orchestrates `postgres`, `backend`, and `frontend` together:

```bash
docker compose up --build
```

This directory is reserved for additional container assets (e.g. database init scripts or extra services) as the project grows.
