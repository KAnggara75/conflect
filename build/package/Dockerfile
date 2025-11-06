# ===========================
# Stage 1: Build
# ===========================
FROM golang:alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary (tanpa CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -ldflags="-s -w" -o conflect ./cmd/conflect


# ===========================
# Stage 2: Runtime
# ===========================
FROM alpine:3.22

# Install git dan CA certificates (wajib untuk HTTPS)
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Tambahkan user non-root
RUN adduser -D -u 1001 appuser

# Set working directory
WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/conflect .

# Buat direktori kerja untuk clone repo dan berikan akses penuh
RUN mkdir -p /app/tmp/repo && chown -R appuser:appuser /app

# Gunakan user non-root
USER appuser

# Environment variables default (bisa dioverride di docker run)
ENV APP_PORT=8080
EXPOSE 8080

# Labels (metadata image)
LABEL org.opencontainers.image.title="Conflect" \
      org.opencontainers.image.description="Spring Cloud Config alternative in Go" \
      org.opencontainers.image.source="https://github.com/KAnggara75/conflect" \
      org.opencontainers.image.licenses="Apache-2.0"

# Jalankan aplikasi
ENTRYPOINT ["./conflect"]
