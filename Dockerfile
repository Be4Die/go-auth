FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/auth-service

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /auth-service .
RUN adduser -D -u 10001 appuser
USER appuser
ENV GIN_MODE=release
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["./auth-service"]
