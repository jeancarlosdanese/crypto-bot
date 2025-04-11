.PHONY: run test build up down clean logs restart check-env cover

# Carrega variáveis do .env local se existir
ifneq (,$(wildcard .env))
	include .env
	export
endif

run: check-env
	go run ./cmd/api/main.go

test:
	go test ./... -v

build: check-env
	CGO_ENABLED=0 GOOS=linux go build -o build/crypto-bot ./cmd/api

up:
	cd .. && docker compose up -d --build

down:
	cd .. && docker compose down

restart:
	cd .. && docker compose down && docker compose up -d --build

logs:
	cd .. && docker compose logs -f --tail=100

clean:
	rm -rf build coverage.out

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

check-env:
	@echo "🔍 Verificando variáveis de ambiente obrigatórias..."
	@if [ -z "$(POSTGRES_USER)" ] || [ -z "$(POSTGRES_PASSWORD)" ] || [ -z "$(POSTGRES_DATABASE)" ]; then \
		echo "❌ Variáveis obrigatórias ausentes. Verifique seu .env."; \
		exit 1; \
	fi
