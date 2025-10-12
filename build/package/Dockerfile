# Stage 1: Build
FROM golang:alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary (tanpa CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o conflect ./cmd/conflect

# Stage 2: Runtime
FROM alpine:3.22

# Install git (required for cloning repos)
RUN apk add --no-cache git

# Tambahkan user non-root
RUN adduser -D appuser

# Set working directory
WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/conflect .

# ✅ Buat direktori tmp dan berikan akses penuh untuk user appuser
RUN mkdir -p /app/tmp/repo && chown -R appuser:appuser /app

# Ganti user ke non-root
USER appuser

# Environment variables
ENV APP_PORT=8080
EXPOSE 8080

# Labels (metadata image)
LABEL org.opencontainers.image.title="Conflect" \
      org.opencontainers.image.description="Spring Cloud Config alternative in Go" \
      org.opencontainers.image.source="https://github.com/KAnggara75/conflect" \
      org.opencontainers.image.licenses="Apache-2.0"

# Jalankan aplikasi
ENTRYPOINT ["./conflect"]
