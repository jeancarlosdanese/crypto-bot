# -------------------------------------
# 🤖 Projeto: Robô de Crypto
# 📍 Makefile Unificado (Raiz do Projeto)
# -------------------------------------

.PHONY: all init start stop dev restart clean tidy test build logs help

# Diretórios
BACKEND_DIR := backend
FRONTEND_DIR := frontend

# ----------------------
# 🔰 Comandos principais
# ----------------------

# Inicia backend + frontend em modo desenvolvimento
dev: backend-up frontend-dev

# Inicializa tudo (build + start)
init: build start

# Inicia backend (docker) + frontend (npm run dev)
start: backend-up frontend-dev

# Para tudo (docker down)
stop: backend-down

# Reinicia containers do backend
restart: backend-restart

# Build geral (backend + frontend)
build: backend-build frontend-build

# Limpa artefatos
clean: backend-clean frontend-clean

# ----------------------
# 🧠 Backend (Go)
# ----------------------

backend-up:
	$(MAKE) -C $(BACKEND_DIR) up

backend-down:
	$(MAKE) -C $(BACKEND_DIR) down

backend-restart:
	$(MAKE) -C $(BACKEND_DIR) restart

backend-logs:
	$(MAKE) -C $(BACKEND_DIR) logs

backend-build:
	$(MAKE) -C $(BACKEND_DIR) build

backend-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-dev:
	$(MAKE) -C $(BACKEND_DIR) dev

backend-test:
	$(MAKE) -C $(BACKEND_DIR) test

backend-tidy:
	$(MAKE) -C $(BACKEND_DIR) tidy

backend-cover:
	$(MAKE) -C $(BACKEND_DIR) cover

backend-clean:
	$(MAKE) -C $(BACKEND_DIR) clean

# ----------------------
# 🌐 Frontend (Next.js)
# ----------------------

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-clean:
	cd $(FRONTEND_DIR) && rm -rf .next node_modules

frontend-install:
	cd $(FRONTEND_DIR) && npm install

# ----------------------
# 🛠️ Extras
# ----------------------

tidy: backend-tidy
test: backend-test
logs: backend-logs

# ----------------------
# 🆘 Ajuda
# ----------------------

help:
	@echo "🤖 Comandos disponíveis:"
	@echo ""
	@echo "🟢 Início:"
	@echo "  make dev                 - Inicia backend + frontend (modo dev)"
	@echo "  make init                - Builda tudo e inicia"
	@echo "  make start               - Inicia containers + frontend dev"
	@echo "  make stop                - Para containers"
	@echo ""
	@echo "🧠 Backend:"
	@echo "  make backend-run         - Roda o backend direto (sem docker)"
	@echo "  make backend-dev         - Roda com Air (hot reload)"
	@echo "  make backend-build       - Compila binário do backend"
	@echo "  make backend-test        - Executa testes do backend"
	@echo "  make backend-tidy        - Roda go mod tidy"
	@echo "  make backend-cover       - Geração de cobertura dos testes"
	@echo "  make backend-clean       - Limpa build e cobertura"
	@echo "  make backend-up          - Sobe containers via Docker"
	@echo "  make backend-down        - Derruba containers Docker"
	@echo "  make backend-restart     - Reinicia containers Docker"
	@echo "  make backend-logs        - Mostra logs do Docker backend"
	@echo ""
	@echo "🌐 Frontend:"
	@echo "  make frontend-dev        - Inicia frontend em modo dev"
	@echo "  make frontend-build      - Compila frontend (Next.js)"
	@echo "  make frontend-clean      - Limpa .next e node_modules"
	@echo "  make frontend-install    - Instala dependências do frontend"
	@echo ""
	@echo "🧹 Manutenção:"
	@echo "  make tidy                - Alias para backend-tidy"
	@echo "  make test                - Alias para backend-test"
	@echo "  make logs                - Alias para backend-logs"
	@echo "  make clean               - Limpa frontend + backend"
