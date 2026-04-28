# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Run stage
FROM alpine:3.21

RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Railway sets PORT env var, app will use it
EXPOSE 8080

CMD ["./main"]
