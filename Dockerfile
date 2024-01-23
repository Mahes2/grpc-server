# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . /app

RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o grpc-server ./cmd/main.go

# Expose the port that the application will run on
EXPOSE 1531

# Run the executable
CMD ["/app/grpc-server"]
