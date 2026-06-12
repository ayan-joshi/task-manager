# ---- Build stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependencies first for faster rebuilds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a fully static binary (CGO disabled; all drivers are pure Go).
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.20

# ca-certificates for outbound TLS; create a non-root user.
RUN apk add --no-cache ca-certificates wget && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /app/server /app/server

# Persisted upload directory.
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app

ENV PORT=8080 UPLOAD_DIR=/app/uploads
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/server"]
