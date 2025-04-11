.PHONY: run test build up down clean logs restart

run:
	go run ./cmd/api/main.go

test:
	go test ./... -v

build:
	CGO_ENABLED=0 GOOS=linux go build -o build/crypto-bot ./cmd/api

up:
	docker compose up -d --build

down:
	docker compose down

restart:
	docker compose down && docker compose up -d --build

logs:
	docker compose logs -f --tail=100

clean:
	rm -rf build coverage.out

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
