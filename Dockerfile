# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy data to working dir
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app .

# Command to run the executable
CMD ["./main"]