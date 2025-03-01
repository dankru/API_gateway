# Stage 1: Build the Go application
FROM golang:1.23.4 AS builder

# Set the working directory for the container
WORKDIR /gateway

# Install air for live reloading (optional, remove if not using)
RUN go install github.com/air-verse/air@latest

# Copy only the go.mod and go.sum for dependency caching
#COPY go.mod /gateway

# Download dependencies
# RUN go mod download

# Copy the rest of the gatewaylication
#COPY . .

WORKDIR /gateway/cmd

# Build the Go application to a binary (optional if not using `air`)
#RUN go build -o /app/main ./cmd/main.go

# Final CMD for live reload using `air` (optional)
CMD ["air"]
