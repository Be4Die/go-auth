SHELL := /bin/sh

.PHONY: run build test cover docker-build docker-up docker-down

run:
	go run ./cmd/auth-service

build:
	CGO_ENABLED=0 go build -o bin/auth-service ./cmd/auth-service

test:
	go test ./... -race -cover

cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker-compose build auth-service

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
