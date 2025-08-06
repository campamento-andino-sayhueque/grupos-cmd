# --- Build Stage ---
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/grupos-cmd ./cmd/server

# --- Final Stage ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/grupos-cmd .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./grupos-cmd"]
