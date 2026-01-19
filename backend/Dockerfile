# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies separately for better caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY backend/ ./

# Build the application
RUN go build -o main cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Set environment variables (defaults)
ENV PORT=8080

# Expose the API port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
