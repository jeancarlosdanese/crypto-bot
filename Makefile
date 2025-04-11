.PHONY: backend-up backend-down backend-logs backend-restart backend-build frontend-dev frontend-build

# 📦 Backend (API Go)
backend-up:
	cd backend && make up

backend-down:
	cd backend && make down

backend-logs:
	cd backend && make logs

backend-restart:
	cd backend && make restart

backend-build:
	cd backend && make build

# 🌐 Frontend (Next.js)
frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build
