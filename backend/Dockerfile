# Start from the official Go image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o expense-tracker

# Expose the application port (Railway will set the actual port via the PORT environment variable)
EXPOSE 8080

# Command to run the executable
CMD ["./expense-tracker"]