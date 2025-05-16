# Build stage
FROM golang:1.23.4-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# Final stage
FROM alpine:latest

# Update and install ca-certificates for HTTPS
RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main /root/main

# Copy environment file if needed
COPY --from=builder /app/.env /root/.env

# Expose the port your application runs on
EXPOSE 8081

# Command to run
CMD ["./main"]
