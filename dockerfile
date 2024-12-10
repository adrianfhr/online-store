# Step 1: Use an official Golang image for building
FROM golang:1.20 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Use a minimal image for the final container
FROM gcr.io/distroless/base-debian12

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 for Cloud Run
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
