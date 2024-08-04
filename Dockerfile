# Start from the official Golang image
FROM golang:1.20 AS builder

# Set environment variables
ENV GO111MODULE=on
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=linux go build -o app ./cmd

# Use a smaller base image to run the application
FROM alpine:latest  

# Install necessary packages for the application to run
RUN apk --no-cache add ca-certificates

# Set the working directory to /root/
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ./app
