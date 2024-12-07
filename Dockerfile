# Use the official Go image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o nzbmonitor ./internal

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/nzbmonitor .

# Copy the provider config from the builder stage
COPY --from=builder /app/provider.json .

# Expose the port the app runs on
EXPOSE 6666

# Command to run the executable
CMD ["./nzbmonitor"]
