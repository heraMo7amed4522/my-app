# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy the entire project to have access to both services
COPY . .
# Set working directory to chat-services
WORKDIR /app/services/chat-services

# Download dependencies
RUN go mod download


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-service .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/services/chat-services/chat-service .

# Copy .env file
COPY --from=builder /app/services/chat-services/.env .

# Expose port
EXPOSE 50061

# Run the application
CMD ["./chat-service"]
