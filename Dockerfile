FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main ./server/main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/main .

# Expose port  
EXPOSE 9000

# Set the entrypoint to the binaryEnsure executable permission
RUN chmod +x /app/main

# Run the binary
CMD ["./main"]