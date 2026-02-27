# ==========================================
# Stage 1: Build the Go Binary
# ==========================================
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the raw source code
COPY main.go .

# Build the binary. 
# CGO_ENABLED=0 ensures it is a statically linked binary (maximum security/portability)
RUN CGO_ENABLED=0 GOOS=linux go build -o agentauth-core main.go

# ==========================================
# Stage 2: Create the Minimal Production Image
# ==========================================
FROM alpine:latest

# Add ca-certificates just in case we need to make secure outbound HTTPS calls later
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy ONLY the pre-built binary from the builder stage
COPY --from=builder /app/agentauth-core .

# Expose the port your Go server listens on
EXPOSE 8080

# Command to run the executable
CMD ["./agentauth-core"]