# Start from the official Go image
FROM golang:1.24 AS builder

# Set the working directory
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN go build -o server ./main.go

# Final stage: create minimal image
FROM debian:bookworm-slim

# Copy binary from builder stage
COPY --from=builder /app/server /server

# Expose API port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/server"]
