# Use specific versions for base images
FROM golang:1.21.9-alpine3.19 AS builder

WORKDIR /app

# Copy only necessary files for module downloading
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application with optimized flags
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/main main.go

# Use a smaller base image for the final stage
FROM alpine:3.19

WORKDIR /app

# Copy built binary and configuration file
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Set the entry point with necessary parameters
ENTRYPOINT ["./main"]