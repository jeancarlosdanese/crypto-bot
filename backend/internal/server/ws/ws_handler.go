// internal/server/ws/ws_handler.go

package ws

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jeancarlosdanese/crypto-bot/internal/auth"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Recebe botID via URL e token via query string
func SecureWebSocketHandler(botRepo repository.BotRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Conectando ao WebSocket...", "url", r.URL.Path)

		// Extrair botID da URL: /ws/{botID}
		botIDStr := r.PathValue("botID")
		botID, err := uuid.Parse(botIDStr)
		if err != nil {
			logger.Error("Erro ao converter botID para UUID:", err)
			http.Error(w, "botID inválido", http.StatusBadRequest)
			return
		}

		// Extrair token da query string
		token := r.URL.Query().Get("token")
		if token == "" {
			logger.Error("Token ausente na requisição", nil)
			http.Error(w, "token ausente", http.StatusUnauthorized)
			return
		}

		// Validar token JWT
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			logger.Error("Erro ao validar token JWT:", err)
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}

		accountID, err := uuid.Parse(claims.AccountID)
		if err != nil {
			logger.Error("Erro ao converter accountID para UUID:", err)
			http.Error(w, "account_id inválido no token", http.StatusUnauthorized)
			return
		}
		// Verificar se o bot pertence à conta
		bot, err := botRepo.GetByID(botID)
		if err != nil || bot.AccountID != accountID {
			logger.Error("Erro ao buscar bot ou bot não pertence à conta:", err)
			http.Error(w, "bot não pertence à sua conta", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Erro ao fazer upgrade para WebSocket:", err)
			return
		}

		logger.Debug("🧩 Cliente conectado via WebSocket", "bot_id", botID.String(), "account_id", accountID)

		AddClient(bot.ID.String(), conn)
	}
}
