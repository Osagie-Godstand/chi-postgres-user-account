# Use an official Golang runtime as a parent image
FROM golang:1.20.3-alpine

# Set the working directory to /app
WORKDIR /api

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o bin/api ./cmd/api

# Expose port 9090
EXPOSE 9090

# Set the entry point of the container to the executable
CMD ["./bin/api"]
