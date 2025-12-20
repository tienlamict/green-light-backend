# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Build seed tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./scripts/seed.go

# Final stage
FROM alpine:latest

# Install ca-certificates, mysql-client, and netcat for HTTPS and migrations
RUN apk --no-cache add ca-certificates mysql-client netcat-openbsd bash

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/seed ./scripts/

# Copy migrations SQL files
COPY --from=builder /app/migrations/*.sql ./migrations/

# Copy source files needed for seed (if running with go run)
COPY --from=builder /app/scripts ./scripts/
COPY --from=builder /app/internal ./internal/
COPY --from=builder /app/pkg ./pkg/

# Copy init script from builder stage
COPY --from=builder /app/docker-init.sh /root/docker-init.sh
# Fix line endings (remove CR) and make executable
RUN sed -i 's/\r$//' /root/docker-init.sh && \
    chmod +x /root/docker-init.sh && \
    ls -la /root/docker-init.sh && \
    test -f /root/docker-init.sh && echo "✅ docker-init.sh copied successfully"

# Create uploads directory
RUN mkdir -p /root/uploads

# Expose port
EXPOSE 8080

# Run init script which handles migration + app startup
CMD ["/root/docker-init.sh"]

