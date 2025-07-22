FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fertit ./cmd/server

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Install CA certificates for SMTP connections
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/fertit .

# Copy web templates and static files
COPY --from=builder /app/web /app/web

# Create a non-root user to run the application
RUN adduser -D -g '' appuser
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/fertit"]