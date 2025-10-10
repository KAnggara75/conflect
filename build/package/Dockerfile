# Stage 1: Builder
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/conflect ./cmd/conflect

# Stage 2: Runtime
FROM scratch
COPY --from=builder /bin/conflect /conflect
ENV APP_PORT=8080
ENTRYPOINT ["/conflect"]
