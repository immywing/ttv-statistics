# Start from official Go image
FROM golang:1.21-alpine

# Set necessary Go env vars
ENV CGO_ENABLED=0 \
    GO111MODULE=on

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum first to cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of app
COPY . .

# Build the Go binary
RUN go build -o ttv-statistics .

# Default port (informational only)
EXPOSE 8080

# Use ENTRYPOINT so command/flags can be overridden easily in Docker Compose
ENTRYPOINT ["/app/ttv-statistics"]
