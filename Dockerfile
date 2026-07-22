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

# Create a non-privileged system group and user
RUN addgroup -S -g 10001 titan && \
    adduser -S -u 10001 -G titan -h /app -s /bin/false titan

WORKDIR /app

# Copy the compiled binary from the builder stage and set ownership
COPY --from=builder --chown=titan:titan /app/bin/titanhttp .

# Create a public directory for static file serving and set ownership
RUN mkdir -p public && chown -R titan:titan /app

# Switch to the non-privileged user
USER titan:titan

# Expose standard HTTP and HTTPS ports
EXPOSE 8080 8443

# Start the server
ENTRYPOINT ["./titanhttp"]
