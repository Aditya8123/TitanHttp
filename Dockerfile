# Builder Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build statically-linked, optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/titanhttp ./cmd/titanhttp

# Production Stage
FROM alpine:3.20

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/titanhttp .

# Create a public directory for static file serving, if any
RUN mkdir -p public

# Expose standard HTTP and HTTPS ports
EXPOSE 8080 8443

# Start the server
ENTRYPOINT ["./titanhttp"]
