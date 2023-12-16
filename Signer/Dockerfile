# Use the official Go base image for the build
FROM golang:1.21.5 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project source code
COPY . .

# Build the application. Make sure to adjust the executable name
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the executable from the build image
COPY --from=builder /app/main .

# Expose the port on which your application will run
EXPOSE 8080

# Command to run the application
CMD ["./main"]
