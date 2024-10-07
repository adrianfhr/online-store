# Use the official Golang image as the base image
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# Expose the port the app runs on
EXPOSE 8081

# Command to run the executable
CMD ["./myapp"]
