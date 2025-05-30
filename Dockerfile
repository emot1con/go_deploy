# Use the official Go image as base
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Use a smaller image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files if they exist
COPY --from=builder /app/static ./static

# Expose port (Railway will set the PORT environment variable)
EXPOSE 8080

# Run the application
CMD ["./main"]
