# Frontend build stage
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy built frontend assets from the node stage
COPY --from=frontend-builder /app/web/dist ./internal/api/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /rms-server ./cmd/server

# Final stage
FROM alpine:3.19
WORKDIR /root/
RUN apk add --no-cache ca-certificates curl
COPY --from=builder /rms-server .
EXPOSE 8080
CMD ["./rms-server"]
