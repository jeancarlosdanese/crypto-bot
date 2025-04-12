// internal/server/routes/bot_routes.go

package routes

import (
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/handlers"
)

// RegisterBotRoutes adiciona as rotas relacionadas aos bots
func RegisterBotRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.Handler) http.HandlerFunc,
	botRepo repository.BotRepository,
) {
	handler := handlers.NewBotHandle(botRepo)

	mux.Handle("GET /bots", authMiddleware(http.HandlerFunc(handler.ListBotsHandle())))
	// Futuro:
	// mux.Handle("POST /bots", authMiddleware(http.HandlerFunc(handler.CreateBotHandler())))
	// mux.Handle("GET /bots/{id}", authMiddleware(http.HandlerFunc(handler.GetBotHandler())))
}
