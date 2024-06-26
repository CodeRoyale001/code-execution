# Use an official Golang runtime as the base image
FROM jals413/goimage

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# COPY .env .env

# Copy the rest of the project files
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080


# Set the command to run the executable
CMD ["./main"]
