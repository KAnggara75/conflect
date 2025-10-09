FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY ../.. .
RUN go build -o conflect ./cmd/app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/conflect .
# mount secrets/config via file jika diperlukan
ENV APP_PORT=8080
CMD ["./conflect"]
