// internal/server/routes/routes.go

package routes

import (
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/middlewares"
)

// NewRouter cria e retorna um roteador HTTP configurado.
func NewRouter(
	otpRepo repository.AccountOTPRepository,
	accountRepo repository.AccountRepository,
	botRepo repository.BotRepository,
) *http.ServeMux {
	mux := http.NewServeMux()

	// 🔥 Criar middlewares
	authMiddleware := middlewares.AuthMiddleware(accountRepo)

	// 🔥 Registrar rotas principais
	RegisterAuthRoutes(mux, authMiddleware, otpRepo)
	RegisterAccountRoutes(mux, authMiddleware, accountRepo)
	RegisterBotRoutes(mux, authMiddleware, botRepo)
	RegisterWebSocketRoutes(mux, botRepo)

	// 🔥 Rota de Health Check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
