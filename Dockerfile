FROM golang:alpine

WORKDIR /app

# Install air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Copy only the go.mod and go.sum files to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Set the working directory to the cmd directory
WORKDIR /app/cmd

# Expose the port
EXPOSE 8080

# Run the application using air
CMD ["air", "-c", "../.air.toml"]
