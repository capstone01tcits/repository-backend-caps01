# ============================================================================
# Build Stage - Compile Go application
# ============================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
# CGO_ENABLED=0: Disable CGO for Alpine compatibility
# GOOS=linux: Ensure Linux binary (cross-compile support)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# ============================================================================
# Runtime Stage - Minimal production image
# ============================================================================
FROM alpine:3.21

# Install runtime dependencies
RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates tzdata wget && \
    rm -rf /var/cache/apk/*

# Set timezone
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/main .

# Documentation:
# - Railway sets PORT env var → app reads it in config.go
# - If PORT not set, defaults to APP_PORT (5000)
# - Both environment variables work for flexibility
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Start application
# Railway will override PORT env var
CMD ["./main"]
