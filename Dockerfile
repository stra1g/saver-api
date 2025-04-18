# First stage: build the application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application with proper flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o saver-api ./cmd/

# Second stage: runtime image
FROM alpine:3.21

# Add certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/saver-api /saver-api

# Run the application
ENTRYPOINT ["/saver-api"]
