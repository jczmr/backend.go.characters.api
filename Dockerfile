# Use the official Golang image as a base
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-dragonball-service ./cmd/api

# Use a minimal image for the final stage
FROM golang:1.22-alpine

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /go-dragonball-service .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./go-dragonball-service"]