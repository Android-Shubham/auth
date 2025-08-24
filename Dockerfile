# Start from the official Go image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
# Copy go mod and sum files
COPY go.mod go.sum ./
# If you have vendor folder, copy it too
COPY vendor ./vendor

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o auth .

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/auth .

# Copy .env if needed
COPY .env .env

EXPOSE 8080

CMD ["./auth"]