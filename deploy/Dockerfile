# ==========================================
# Stage 1: The Builder
# ==========================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy the go.mod and go.sum files first 
# (This caches dependencies so future builds are faster)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your Go source code
COPY . .

# Compile the Go application
# CGO_ENABLED=0 ensures we build a static binary (no external C libraries required)
RUN CGO_ENABLED=0 GOOS=linux go build -o k8bot ./cmd/k8bot/main.go

# ==========================================
# Stage 2: The Final Image
# ==========================================
FROM alpine:latest

# Install CA Certificates
# (CRITICAL: Your bot needs this to securely talk to the Slack HTTPS API)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy ONLY the compiled binary from Stage 1
COPY --from=builder /app/k8bot .

# Command to run the bot when the container starts
CMD ["./k8bot"]