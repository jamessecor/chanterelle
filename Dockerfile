# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the server binary
COPY --from=builder /app/server .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./server"]