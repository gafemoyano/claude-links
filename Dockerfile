# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy public directory for static files and .well-known
COPY --from=builder /app/public ./public

# Copy docs directory for Swagger documentation
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./main"]