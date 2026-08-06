# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install ca-certificates and git if needed
RUN apk add --no-cache git ca-certificates

# Copy dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/matrix-service ./cmd/server

# Stage 2: Minimal runtime image
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/matrix-service .
COPY --from=builder /app/.env.example .env

# Expose default port
EXPOSE 12000

# Set environment defaults
ENV PORT=12000
ENV NODE_API_URL=http://localhost:4000
ENV HTTP_TIMEOUT_MS=5000
ENV AUTH_ENABLED=false

ENTRYPOINT ["./matrix-service"]
