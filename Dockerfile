# Build stage
FROM golang:1.23.0-alpine3.19 AS builder

WORKDIR /app

# Set Go env variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download \
    && go mod verify

# Copy source code
COPY . .

# Build the application
RUN go build -ldflags='-w -s' -a -o order-management-ms ./src/main

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/order-management-ms .

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./order-management-ms"]
