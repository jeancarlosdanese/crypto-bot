.PHONY: run dev test build up down clean logs restart check-env cover tidy help

# E define no início do Makefile
DB_SERVICE := postgres

# Carrega variáveis do .env local se existir
ifneq (,$(wildcard .env))
	include .env
	export
endif

# ----------------------
# 🔧 Comandos principais
# ----------------------

run: check-env
	go run ./cmd/api/main.go

dev: up-only
	air

up-only:
	docker compose up -d $(DB_SERVICE)

test:
	go test ./... -v

build: check-env
	CGO_ENABLED=0 GOOS=linux go build -o build/crypto-bot ./cmd/api

tidy:
	go mod tidy

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf build coverage.out

check-env:
	@echo "🔍 Verificando variáveis de ambiente obrigatórias..."
	@if [ -z "$(POSTGRES_USER)" ] || [ -z "$(POSTGRES_PASSWORD)" ] || [ -z "$(POSTGRES_DATABASE)" ]; then \
		echo "❌ Variáveis obrigatórias ausentes. Verifique seu .env."; \
		exit 1; \
	fi

# ----------------------
# 🐳 Docker
# ----------------------

up:
	cd .. && docker compose up -d --build

down:
	cd .. && docker compose down

restart:
	cd .. && docker compose down && docker compose up -d --build

logs:
	cd .. && docker compose logs -f --tail=100

# ----------------------
# 🆘 Ajuda
# ----------------------

help:
	@echo "🧠 Comandos do Backend:"
	@echo ""
	@echo "🟢 Execução:"
	@echo "  make run           - Roda a aplicação (modo local)"
	@echo "  make dev           - Hot reload com Air"
	@echo "  make build         - Compila binário para produção"
	@echo ""
	@echo "🧪 Testes:"
	@echo "  make test          - Executa testes"
	@echo "  make cover         - Geração de cobertura de testes"
	@echo ""
	@echo "🧹 Manutenção:"
	@echo "  make tidy          - go mod tidy"
	@echo "  make clean         - Limpa build e cobertura"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make up            - Sobe containers via Docker Compose"
	@echo "  make down          - Derruba containers"
	@echo "  make restart       - Reinicia containers"
	@echo "  make logs          - Logs do backend (docker)"
