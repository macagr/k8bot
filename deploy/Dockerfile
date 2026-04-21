# ==========================================
# Stage 1: The Builder
# ==========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# 🧪 DEBUG STEP: We copy everything, then list it.
# If /app/cmd does not appear in your logs, the context is wrong.
COPY . .
RUN ls -R /app

# Compile the Go application
# We use the relative path starting from /app
RUN CGO_ENABLED=0 GOOS=linux go build -o k8bot ./cmd/k8bot/main.go

# ==========================================
# Stage 2: The Final Image
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/k8bot .

CMD ["./k8bot"]