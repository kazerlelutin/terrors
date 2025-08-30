# 🎭 Terrors - Error Monitoring Service Dockerfile
# "The call is coming from inside the container..."

# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Update go.mod and build the application
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S terrors && \
    adduser -u 1001 -S terrors -G terrors

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy static files
COPY --from=builder /app/static ./static

# Create empty .env file to prevent warnings (variables will come from CapRover)
RUN touch .env

# Change ownership to non-root user
RUN chown -R terrors:terrors /app

# Switch to non-root user
USER terrors

# Expose port (CapRover will use this)
EXPOSE 3000

# Health check for CapRover
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Run the application
CMD ["./main"]
