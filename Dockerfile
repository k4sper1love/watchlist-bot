# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files (go.mod and go.sum) to the working directory
COPY go.mod .
COPY go.sum .

# Download the Go dependencies specified in go.mod and go.sum
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Build the Go application with CGO disabled (for static linking) and target OS as Linux
# The output binary is named "watchlist-bot" and is built from the ./cmd/bot directory
RUN CGO_ENABLED=0 GOOS=linux go build -o watchlist-bot ./cmd/bot

# Stage 2: Create a minimal runtime environment using Alpine Linux
FROM alpine:latest

# Install CA certificates to ensure secure HTTPS connections
RUN apk --no-cache add ca-certificates

# Set the working directory inside the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage to the runtime container
COPY --from=builder /app/watchlist-bot .
COPY --from=builder /app/locales ./locales

# Define the default command to run the application when the container starts
CMD ["./watchlist-bot"]