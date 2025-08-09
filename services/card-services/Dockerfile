# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy the entire project to have access to both services
COPY . .
# Set working directory to card-services
WORKDIR /app/services/card-services

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o card-service .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/services/card-services/card-service .

# Copy .env file
COPY --from=builder /app/services/card-services/.env .

# Expose port
EXPOSE 50053

# Run the application
CMD ["./card-service"]