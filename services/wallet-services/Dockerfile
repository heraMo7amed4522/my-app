FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project to have access to both services
COPY . .

# Set working directory to wallet-services
WORKDIR /app/services/wallet-services

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-service .

# Use a minimal base image for the final stage
FROM alpine:latest


RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/wallet-services/wallet-service .

# Copy .env file
COPY --from=builder /app/services/wallet-services/.env .

# Expose the port the app runs on
EXPOSE 50055

# Command to run the application
CMD ["./wallet-service"]