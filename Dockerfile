# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install SQLite runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port (Render will override with PORT env var)
EXPOSE 8080

# Run the server
CMD ["./server"]
