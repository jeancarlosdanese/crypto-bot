// internal/server/handlers/bot_handler.go

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/dto"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/runtime"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/middlewares"
	"github.com/jeancarlosdanese/crypto-bot/internal/utils"
)

type BotHandle interface {
	ListBotsHandle() http.HandlerFunc
	GetCandlesHandler() http.HandlerFunc
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

func (h *botHandle) GetCandlesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := utils.GetUUIDFromRequestPath(r, w, "id")

		runtime.BotsMap.RLock()
		strategy := runtime.BotsMap.Items[id]
		runtime.BotsMap.RUnlock()
		if strategy == nil {
			utils.SendError(w, http.StatusNotFound, "Bot não está ativo")
			return
		}

		var result []map[string]interface{}
		skipped := 0
		for _, c := range strategy.CandlesWindow {
			if c.Time == 0 {
				skipped++
				continue // ignora candles sem timestamp
			}
			result = append(result, map[string]interface{}{
				"time":  c.Time,
				"open":  c.Open,
				"high":  c.High,
				"low":   c.Low,
				"close": c.Close,
			})
		}
		if skipped > 0 {
			logger.Debug("⏩ Ignorados candles sem timestamp", "qtd", skipped)
		}

		utils.SendJSON(w, http.StatusOK, result)
	}
}
