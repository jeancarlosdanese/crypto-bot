// internal/server/handlers/bot_handler.go

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/dto"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/middlewares"
)

type BotHandle interface {
	ListBotsHandle() http.HandlerFunc
}

type botHandle struct {
	repo repository.BotRepository
}

func NewBotHandle(repo repository.BotRepository) BotHandle {
	return &botHandle{repo: repo}
}

func (h *botHandle) ListBotsHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, ok := middlewares.GetAuthenticatedAccount(r.Context())
		if !ok {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		bots, err := h.repo.GetByAccountID(account.ID)
		if err != nil {
			logger.Error("Erro ao buscar bots", err)
			http.Error(w, "Erro ao buscar bots", http.StatusInternalServerError)
			return
		}

		var result []dto.BotResponseDTO
		for _, bot := range bots {
			result = append(result, dto.NewBotResponseDTO(&bot))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
