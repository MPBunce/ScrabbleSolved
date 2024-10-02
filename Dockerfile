# Use Ubuntu as the base image
FROM ubuntu:latest

# Install Go and CA certificates in the Ubuntu image
RUN apt-get update && apt-get install -y golang ca-certificates

# Set the working directory
WORKDIR /app

# Copy the Go modules and dependencies
COPY go.mod ./
COPY go.sum ./

# Install the certificates required for secure SSL connections
RUN update-ca-certificates

# Download Go modules
RUN go mod download

# Copy the rest of the source code
COPY . .

ENV PORT=4000

# Expose the application port
EXPOSE 4000

# Build the Go app before running it
RUN go build -o main cmd/app/*

# Run the Go app binary
CMD ["./main"]
