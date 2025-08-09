# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . .

# Set working directory to transaction-services
WORKDIR /app/services/transaction-services

# Download all dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o transaction-service .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Set the working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/transaction-services/transaction-service .

# Copy .env file
COPY --from=builder /app/services/transaction-services/.env .

# Expose the port the app runs on
EXPOSE 50054

# Command to run the executable
CMD ["./transaction-service"]