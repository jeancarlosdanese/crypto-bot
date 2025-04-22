// internal/server/handlers/bot_handler.go

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/dto"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/runtime"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/middlewares"
	"github.com/jeancarlosdanese/crypto-bot/internal/utils"
)

type BotHandle interface {
	ListBotsHandle() http.HandlerFunc
	GetBotByIDHandle() http.HandlerFunc
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

func (h *botHandle) GetBotByIDHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := utils.GetUUIDFromRequestPath(r, w, "id")

		runtime.BotsMap.RLock()
		strategy := runtime.BotsMap.Items[id]
		runtime.BotsMap.RUnlock()
		if strategy == nil {
			utils.SendError(w, http.StatusNotFound, "Bot não está ativo")
			return
		}

		bot, err := h.repo.GetByID(id)
		if err != nil {
			logger.Error("Erro ao buscar bot", err)
			http.Error(w, "Erro ao buscar bot", http.StatusInternalServerError)
			return
		}
		if bot == nil {
			utils.SendError(w, http.StatusNotFound, "Bot não encontrado")
			return
		}

		result := dto.NewBotResponseDTO(bot)
		utils.SendJSON(w, http.StatusOK, result)
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
		prices := strategy.ClosingPrices()
		for i, c := range strategy.CandlesWindow {
			if c.Time == 0 {
				skipped++
				continue // ignora candles sem timestamp
			}
			candleMap := map[string]interface{}{
				"time":  c.Time,
				"open":  c.Open,
				"high":  c.High,
				"low":   c.Low,
				"close": c.Close,
			}

			// Cálculo das EMAs da estratégia EMA_FAN
			emaPeriods := []int{10, 15, 20, 25, 30, 35, 40}
			for _, period := range emaPeriods {
				if i >= period-1 {
					candleMap[fmt.Sprintf("ema%d", period)] = indicators.MovingAverage(prices[:i+1], period)
				}
			}

			// Cálculo das médias clássicas (se ainda usadas)
			if i >= 8 {
				candleMap["ma9"] = indicators.MovingAverage(prices[:i+1], 9)
			}
			if i >= 25 {
				candleMap["ma26"] = indicators.MovingAverage(prices[:i+1], 26)
			}

			// Cálculo do RSI
			rsi := 0.0
			if i >= 2 {
				rsi = indicators.RSI(prices[:i+1], 2)
			}
			candleMap["rsi"] = rsi

			// Cálculo do MACD
			macd, signal, hist := indicators.MACD(prices[:i+1], 9, 14, 7)
			if len(macd) > 0 {
				candleMap["macd"] = macd[len(macd)-1]
			}
			if len(signal) > 0 {
				candleMap["macd_signal"] = signal[len(signal)-1]
			}
			if len(hist) > 0 {
				candleMap["macd_histogram"] = hist[len(hist)-1]
			}

			result = append(result, candleMap)
		}
		if skipped > 0 {
			logger.Debug("⏩ Ignorados candles sem timestamp", "qtd", skipped)
		}

		utils.SendJSON(w, http.StatusOK, result)
	}
}
