# Use an official Golang runtime as a parent image
FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-assignment

EXPOSE 9000

# Command to run the executable
CMD ["/go-assignment"]
