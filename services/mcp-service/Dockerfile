# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-service ./cmd/main.go

# Runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mcp-service .

# Expose MCP service port
EXPOSE 8081

# Run as non-root user
USER 1001:1001

CMD ["/app/mcp-service"]
