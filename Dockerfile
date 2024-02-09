FROM golang:alpine

WORKDIR /app

# Install air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Copy only the go.mod and go.sum files to leverage Docker caching
COPY go.mod go.sum ./

RUN apk --no-cache update && \
    apk --no-cache add git gcc libc-dev

# Kafka Go client is based on the C library librdkafka
ENV CGO_ENABLED 1
ENV GOFLAGS -mod=vendor
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk add librdkafka-dev pkgconf
# Download dependencies
RUN go mod download



# Copy the entire project
COPY . .


# Copy the entire project
RUN go mod vendor

# Set the working directory to the cmd directory
WORKDIR /app/cmd

# Expose the port
EXPOSE 8080

# Run the application using air
CMD ["air", "-c", "../.air.toml"]
