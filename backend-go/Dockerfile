# =====================================================================
# STAGE 1: High-Performance Compilation
# =====================================================================
FROM golang:1.21-alpine AS builder

# Install SSL certificates and build tools
RUN apk add --no-cache git ca-certificates build-base

WORKDIR /src

# Copy go.mod and main.go first to run go mod tidy inside the container
COPY go.mod ./
COPY main.go ./

# Resolve dependencies and generate go.sum cryptosecure checksums
RUN go mod tidy
RUN go mod download

# Build statically linked secure binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/backend-go .

# =====================================================================
# STAGE 2: Hardened Minimal Runtime Container
# =====================================================================
FROM alpine:3.18 AS final

# Install standard CA certificates for outbound SSL/TLS calls
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Create a non-root system user for container execution safety (least-privilege)
RUN addgroup -S appuser && adduser -S appuser -G appuser
USER appuser

# Copy static binary from compilation layer
COPY --from=builder /app/backend-go /app/backend-go

# Expose microservice endpoint
EXPOSE 8080

# Environment triggers
ENV GIN_MODE=release
ENV PORT=8080

ENTRYPOINT ["/app/backend-go"]
