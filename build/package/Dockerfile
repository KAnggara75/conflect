# ===========================
# Stage 1: Build
# ===========================
FROM golang:1.25.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o conflect ./cmd/conflect


# ===========================
# Stage 2: Runtime
# ===========================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -u 1001 -m -s /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /app/conflect /app/conflect

RUN mkdir -p /app/tmp/repo && chown -R appuser:appuser /app

USER appuser

ENV APP_PORT=8080
ENV OTEL_ENABLED=false
ENV OTEL_EXPORTER_OTLP_ENDPOINT=""
ENV GOMAXPROCS=2

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${APP_PORT}/health || exit 1

LABEL org.opencontainers.image.title="Conflect" \
      org.opencontainers.image.description="Spring Cloud Config alternative in Go" \
      org.opencontainers.image.source="https://github.com/KAnggara75/conflect" \
      org.opencontainers.image.licenses="Apache-2.0"

ENTRYPOINT ["/app/conflect"]
