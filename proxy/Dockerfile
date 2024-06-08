FROM golang:1.22-alpine

# Install git
RUN apk update && apk add --no-cache git

# Install fresh
RUN go install github.com/gravityblast/fresh@latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY .. .

# Expose the application port (adjust according to your application)
EXPOSE 3001

# Start the application using fresh
CMD ["fresh"]
