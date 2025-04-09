# Makefile na raiz do projeto

.PHONY: run test build up down clean

run:
	go run ./cmd/api/main.go

test:
	go test ./... -v

build:
	go build -o build/crypto-bot ./cmd/api

up:
	docker-compose up -d

down:
	docker-compose down

clean:
	rm -rf build

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
