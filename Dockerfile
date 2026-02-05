# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /rms-server ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /root/

# Install runtime dependencies (e.g., for health checks or debugging)
RUN apk add --no-cache ca-certificates curl

# Copy the binary from the builder stage
COPY --from=builder /rms-server .

# Copy migrations (embedded in binary, but keeping for reference if needed)
# COPY internal/db/migrations ./migrations

# Expose API port
EXPOSE 8080

# Run the server
CMD ["./rms-server"]
