# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

# Copy go.mod first to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download all

# Copy source code
COPY . .
# Skip tidy due to restricted networks; rely on downloaded modules

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /auth-service ./cmd/auth-service

# Run stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /auth-service .
# Copy migrations if we decide to run them from app, or use separate migration tool
# COPY migrations ./migrations 

RUN adduser -D -u 10001 appuser
USER appuser
ENV GIN_MODE=release
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./auth-service"]
