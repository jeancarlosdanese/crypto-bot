// internal/server/routes/ws_routes.go

package routes

import (
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/ws"
)

func RegisterWebSocketRoutes(
	mux *http.ServeMux,
	botRepo repository.BotRepository,
) {
	mux.Handle("GET /ws/{botID}", http.HandlerFunc(ws.SecureWebSocketHandler(botRepo)))
}
