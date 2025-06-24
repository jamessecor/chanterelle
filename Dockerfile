FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.19

WORKDIR /app

# Copy the server binary and .env file
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]
