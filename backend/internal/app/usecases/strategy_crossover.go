package usecases

import (
	"fmt"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	reporter "github.com/jeancarlosdanese/crypto-bot/internal/report"
	serverws "github.com/jeancarlosdanese/crypto-bot/internal/server/ws"
)

func (s *StrategyUseCase) EvaluateCrossover(timestamp int64) string {
	prices := s.ClosingPrices()
	if len(prices) < 26 {
		return "HOLD"
	}

	ma9 := indicators.MovingAverage(prices, 9)
	ma26 := indicators.MovingAverage(prices, 26)
	rsi := indicators.RSI(prices, 14)
	volatility := indicators.Volatility(prices)
	atr := indicators.ATRFromCandles(s.CandlesWindow)
	currentPrice := prices[len(prices)-1]

	basicSignal := "HOLD"
	if ma9 > ma26 && currentPrice > ma9 && rsi < 70 {
		basicSignal = "BUY"
	} else if ma9 < ma26 && currentPrice < ma9 {
		basicSignal = "SELL"
	}

	indicatorsMap := map[string]float64{
		"ma9":        ma9,
		"ma26":       ma26,
		"rsi":        rsi,
		"volatility": volatility,
		"atr":        atr,
		"price":      currentPrice,
	}

	params := map[string]any{
		"ma_short":      9,
		"ma_long":       26,
		"rsi_threshold": 70,
	}
	ctx := map[string]any{
		"candles_total": s.TotalCandles,
		"calibrated_at": s.LastCalibrationGlob,
	}

	strategyName := "EvaluateCrossover"
	strategyVersion := "1.0.1"

	if s.PositionQuantity == 0 && basicSignal == "BUY" {
		// 🔧 Parâmetros dinâmicos
		minVolatility := getFloatParam(params, "volatility_min", 0.0)
		minATR := getFloatParam(params, "atr_min", 0.0)

		// ❌ Ignora entrada se volatilidade ou ATR estiverem abaixo do mínimo
		if volatility < minVolatility {
			logger.Debug("🚫 Entrada bloqueada por baixa volatilidade",
				"symbol", s.Bot.Symbol,
				"volatility", volatility,
				"min_required", minVolatility,
			)
			return "HOLD"
		}

		if atr < minATR {
			logger.Debug("🚫 Entrada bloqueada por ATR insuficiente",
				"symbol", s.Bot.Symbol,
				"atr", atr,
				"min_required", minATR,
			)
			return "HOLD"
		}

		// ✅ Entrada aprovada
		s.PositionQuantity = 1
		s.LastEntryPrice = currentPrice
		s.LastEntryTimestamp = timestamp
		s.LastDecision = "BUY"

		logger.Info("📈 Entrada executada (Crossover)",
			"symbol", s.Bot.Symbol,
			"price", currentPrice,
			"ma9", ma9,
			"ma26", ma26,
			"rsi", rsi,
			"volatility", volatility,
			"atr", atr,
		)

		// 💾 Salvar posição aberta
		if s.PositionRepo != nil {
			logger.Debug("🧩 Salvando posição em aberto...",
				"bot_id", s.Bot.ID.String(),
				"price", currentPrice,
				"timestamp", timestamp,
			)

			err := s.PositionRepo.Save(entity.OpenPosition{
				BotID:      s.Bot.ID,
				EntryPrice: currentPrice,
				Timestamp:  timestamp,
			})
			if err != nil {
				logger.Error("❌ Erro ao salvar posição", err, "bot_id", s.Bot.ID.String())
			} else {
				logger.Debug("✅ Posição salva com sucesso", "bot_id", s.Bot.ID.String())
			}
		}

		// 📝 Log de decisão
		s.saveDecisionLog(strategyName, strategyVersion, "BUY", timestamp, indicatorsMap, params, ctx)

		// 💬 Enviar evento de decisão para o WebSocket
		serverws.Publish(s.Bot.ID.String(), serverws.Event{
			Type: "decision",
			Data: map[string]interface{}{
				"time":     timestamp / 1000,
				"price":    currentPrice,
				"decision": "BUY",
			},
		})

		return "BUY"
	}

	if s.PositionQuantity > 0 {
		// 🔧 Parâmetros configuráveis
		emaTrailingPeriod := getIntParam(params, "ema_trailing", 5)
		rsiExitThreshold := getFloatParam(params, "rsi_exit_threshold", 80)
		atrMultiplier := getFloatParam(params, "atr_multiplier", 1.5)

		// 📊 Indicadores auxiliares
		emaTrailing := indicators.MovingAverage(prices, emaTrailingPeriod)
		rsiPrev := indicators.RSI(prices[:len(prices)-1], 14)
		stopLossThreshold := s.LastEntryPrice + atr*atrMultiplier

		// 🧠 Critérios de saída
		stopLossHit := currentPrice < stopLossThreshold
		priceBelowTrailing := currentPrice < emaTrailing
		rsiReversal := rsiPrev > rsiExitThreshold && rsi < rsiPrev

		if stopLossHit || priceBelowTrailing || rsiReversal || basicSignal == "SELL" {
			reason := ""
			switch {
			case stopLossHit:
				reason = fmt.Sprintf("ATR stop hit (%.2f < %.2f)", currentPrice, stopLossThreshold)
			case priceBelowTrailing:
				reason = fmt.Sprintf("Price below EMA%d (%.2f < %.2f)", emaTrailingPeriod, currentPrice, emaTrailing)
			case rsiReversal:
				reason = fmt.Sprintf("RSI reversal (%.2f < %.2f)", rsi, rsiPrev)
			case basicSignal == "SELL":
				reason = "Crossover reversal signal"
			}

			s.PositionQuantity = 0
			s.LastDecision = "SELL"

			logger.Info("📉 Saída executada (Crossover)",
				"symbol", s.Bot.Symbol, "price", currentPrice, "reason", reason,
				"roi", ((currentPrice-s.LastEntryPrice)/s.LastEntryPrice)*100)

			if s.PositionRepo != nil {
				_ = s.PositionRepo.Delete(s.Bot.ID)
			}
			s.saveDecisionLog(strategyName, strategyVersion, "SELL", timestamp, indicatorsMap, params, ctx)

			profit := currentPrice - s.LastEntryPrice
			roi := (profit / s.LastEntryPrice) * 100
			duration := (timestamp - s.LastEntryTimestamp) / 1000

			exec := entity.ExecutionLog{
				BotID:     s.Bot.ID,
				Symbol:    s.Bot.Symbol,
				Interval:  s.Bot.Interval,
				Entry:     entity.TradePoint{Price: s.LastEntryPrice, Timestamp: s.LastEntryTimestamp},
				Exit:      entity.TradePoint{Price: currentPrice, Timestamp: timestamp},
				Profit:    profit,
				ROIPct:    roi,
				Duration:  duration,
				Strategy:  entity.StrategyInfo{Name: strategyName, Version: strategyVersion},
				CreatedAt: time.Now(),
			}
			if s.ExecutionLogRepo != nil {
				_ = s.ExecutionLogRepo.Save(exec)
			}

			// 💬 Enviar evento de decisão para o WebSocket
			serverws.Publish(s.Bot.ID.String(), serverws.Event{
				Type: "decision",
				Data: map[string]interface{}{
					"time":     timestamp / 1000,
					"price":    currentPrice,
					"decision": "SELL",
				},
			})

			go reporter.PrintExecutionSummary(s.ExecutionLogRepo)
			return "SELL"
		}
	}
	return "HOLD"
}

func getFloatParam(params map[string]any, key string, defaultVal float64) float64 {
	if val, ok := params[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return defaultVal
}

func getIntParam(params map[string]any, key string, defaultVal int) int {
	if val, ok := params[key]; ok {
		if f, ok := val.(float64); ok { // JSON unmarshals numbers as float64
			return int(f)
		}
	}
	return defaultVal
}
