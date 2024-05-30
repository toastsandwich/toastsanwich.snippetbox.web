# Use the official Golang image
FROM golang:1.22.2

# Install Air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o exe .

# Expose port 4000 to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["air", "-c", ".air.toml"]
