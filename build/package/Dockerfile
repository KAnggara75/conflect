# Stage 1: Build
FROM golang:alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o conflect ./cmd/conflect

# Stage 2: Runtime
FROM alpine:3.22

# Install git (required for cloning repos)
RUN apk add --no-cache git

# Tambahkan user non-root
RUN adduser -D appuser

WORKDIR /app
COPY --from=builder /app/conflect .

USER appuser

ENV APP_PORT=8080
EXPOSE 8080

LABEL org.opencontainers.image.title="Conflect" \
      org.opencontainers.image.description="Spring Cloud Config alternative in Go" \
      org.opencontainers.image.source="https://github.com/KAnggara75/conflect" \
      org.opencontainers.image.licenses="Apache-2.0"

ENTRYPOINT ["./conflect"]
