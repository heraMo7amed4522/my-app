FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the entire project to have access to both services
COPY . .

# Set working directory to auth-services
WORKDIR /app/services/auth-services

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/auth-services/auth-service .

# Copy Firebase service account key
COPY ./services/auth-services/my-app-57137-firebase-adminsdk-fbsvc-a9ef791802.json ./firebase-credentials.json

EXPOSE 50052

CMD ["./auth-service"]