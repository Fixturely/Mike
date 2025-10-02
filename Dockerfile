# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Development stage
FROM golang:1.25-alpine AS development

WORKDIR /app

# Install make and air for hot reloading
RUN apk add --no-cache make
RUN go install github.com/air-verse/air@latest

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Use air for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
