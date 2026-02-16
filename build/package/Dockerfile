# ===========================
# Stage 1: Build
# ===========================
FROM golang:1.25.1-bookworm AS builder

# Set working directory
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# WAJIB: CGO aktif, JANGAN strip simbol
RUN CGO_ENABLED=1 GOOS=linux go build -o conflect ./cmd/conflect


# ===========================
# Stage 2: Runtime (glibc, bukan Alpine)
# ===========================
FROM debian:bookworm-slim

# CA certs tetap dibutuhkan
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Tambahkan user non-root (tetap sama konsepnya)
RUN useradd -u 1001 -m appuser

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
