package usecases

import (
	"fmt"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/report"
)

func (s *StrategyUseCase) EvaluateEMAFanWithVolume(timestamp int64) string {
	prices := s.closingPrices()
	if len(prices) < 40 || len(s.CandlesWindow) < 11 {
		return "HOLD"
	}

	periods := []int{10, 15, 20, 25, 30, 35, 40}
	emas := make([]float64, len(periods))
	for i, p := range periods {
		emas[i] = indicators.MovingAverage(prices, p)
	}

	isAligned := true
	for i := 1; i < len(emas); i++ {
		if emas[i] <= emas[i-1] {
			isAligned = false
			break
		}
	}

	slopeMap := calculateEMASlopes(prices, periods)
	for _, slope := range slopeMap {
		if slope < 0.08 {
			return "HOLD"
		}
	}

	lastVolume := s.CandlesWindow[len(s.CandlesWindow)-1].Volume
	avgVolume := 0.0
	for i := len(s.CandlesWindow) - 11; i < len(s.CandlesWindow)-1; i++ {
		avgVolume += s.CandlesWindow[i].Volume
	}
	avgVolume /= 10
	volumeConfirmed := lastVolume > avgVolume

	currentPrice := prices[len(prices)-1]

	indicatorsMap := map[string]float64{
		"price":        currentPrice,
		"last_volume":  lastVolume,
		"avg_volume":   avgVolume,
		"volume_ratio": lastVolume / avgVolume,
	}
	for i, p := range periods {
		indicatorsMap[fmt.Sprintf("ema%d", p)] = emas[i]
		indicatorsMap[fmt.Sprintf("slope%d", p)] = slopeMap[p]
	}

	parameters := map[string]any{
		"emas":          periods,
		"volume_period": 10,
		"slope_min":     0.08,
	}

	context := map[string]any{
		"candles_total": s.TotalCandles,
		"calibrated_at": s.LastCalibrationGlob,
	}

	name := "EvaluateEMAFanWithVolume"
	version := "1.0.0"

	if isAligned && volumeConfirmed && s.PositionQuantity == 0 {
		s.PositionQuantity = 1
		s.LastEntryPrice = currentPrice
		s.LastEntryTimestamp = timestamp
		s.LastDecision = "BUY"

		logger.Info("📈 Entrada (EMA Fan)", "symbol", s.Bot.Symbol, "price", currentPrice, "volume_ratio", lastVolume/avgVolume)

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Save(entity.OpenPosition{
				BotID:      s.Bot.ID,
				EntryPrice: currentPrice,
				Timestamp:  timestamp,
			})
		}

		s.saveDecisionLog(name, version, "BUY", timestamp, indicatorsMap, parameters, context)
		return "BUY"
	}

	if s.PositionQuantity > 0 && !isAligned {
		s.PositionQuantity = 0
		s.LastDecision = "SELL"

		logger.Info("📉 Saída (EMA Fan)", "symbol", s.Bot.Symbol, "price", currentPrice)

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Delete(s.Bot.ID)
		}

		s.saveDecisionLog(name, version, "SELL", timestamp, indicatorsMap, parameters, context)

		profit := currentPrice - s.LastEntryPrice
		roi := (profit / s.LastEntryPrice) * 100
		duration := (timestamp - s.LastEntryTimestamp) / 1000

		exec := entity.ExecutionLog{
			BotID:    s.Bot.ID,
			Symbol:   s.Bot.Symbol,
			Interval: s.Bot.Interval,
			Entry: entity.TradePoint{
				Price:     s.LastEntryPrice,
				Timestamp: s.LastEntryTimestamp,
			},
			Exit: entity.TradePoint{
				Price:     currentPrice,
				Timestamp: timestamp,
			},
			Profit:   profit,
			ROIPct:   roi,
			Duration: duration,
			Strategy: entity.StrategyInfo{Name: name, Version: version},
			CreatedAt: time.Now(),
		}

		if s.ExecutionLogRepo != nil {
			_ = s.ExecutionLogRepo.Save(exec)
		}

		go reporter.PrintExecutionSummary(s.ExecutionLogRepo)
		return "SELL"
	}

	return "HOLD"
}