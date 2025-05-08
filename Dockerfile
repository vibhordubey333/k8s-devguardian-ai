FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o devguardian -ldflags="-w -s" ./main.go

# Create a minimal image
FROM alpine:latest

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/devguardian .

# Copy policies directory
COPY --from=builder /app/internal/policies ./internal/policies

# Expose port if needed
# EXPOSE 8080

# Command to run
ENTRYPOINT ["./devguardian"]
CMD ["audit"]