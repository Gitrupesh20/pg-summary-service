# Use Go 1.24 (matches your go.mod)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source
COPY . .

# Download dependencies
RUN go mod download

# Build binary
RUN go build -o main ./cmd/main.go

# Final lightweight image
FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Run the app
CMD ["./main"]
