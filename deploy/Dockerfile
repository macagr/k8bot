# ==========================================
# Stage 1: The Builder
# ==========================================
# --platform=$BUILDPLATFORM makes the compiler run natively on the runner's CPU
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# TARGETARCH is automatically filled by Docker Buildx (amd64 or arm64)
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🚀 The Magic: We tell Go to cross-compile natively. 
# This is 10x faster than QEMU emulation.
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o k8bot ./cmd/k8bot/main.go

# ==========================================
# Stage 2: The Final Image
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/k8bot .

CMD ["./k8bot"]