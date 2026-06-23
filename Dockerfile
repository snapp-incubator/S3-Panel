# 1) Build the frontend (Vite) into /fe/dist.
FROM node:22-alpine AS frontend

# Build-time configuration, exposed to Vite as environment variables.
ARG VITE_API_LANGUAGE=fa
ARG VITE_APP_ENVIRONMENT
ARG VITE_AUTH_TOKEN
ARG VITE_VARIANT=cab
ARG VITE_CENTRAL_BACKEND_API
ARG VITE_ENV
ENV VITE_API_LANGUAGE=$VITE_API_LANGUAGE \
    VITE_APP_ENVIRONMENT=$VITE_APP_ENVIRONMENT \
    VITE_AUTH_TOKEN=$VITE_AUTH_TOKEN \
    VITE_VARIANT=$VITE_VARIANT \
    VITE_CENTRAL_BACKEND_API=$VITE_CENTRAL_BACKEND_API \
    VITE_ENV=$VITE_ENV

RUN corepack enable && apk add --no-cache git
WORKDIR /fe
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm run build

# 2) Build the Go binary with the frontend embedded (internal/web/dist).
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /fe/dist ./internal/web/dist
RUN go build -v -o s3-panel ./cmd/s3-panel

# 3) Minimal runtime image — a single binary that serves the API and the SPA.
FROM alpine:3.20
WORKDIR /app/
COPY --from=builder /app/s3-panel .
ENTRYPOINT ["./s3-panel"]
CMD ["s3-panel", "--configPath=./config.toml"]
