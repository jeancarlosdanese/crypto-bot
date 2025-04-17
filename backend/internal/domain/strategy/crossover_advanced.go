// internal/domain/strategy/crossover_advanced.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type CrossoverStrategyAdvanced struct{}

func (s *CrossoverStrategyAdvanced) Name() string {
	return "CROSSOVER_ADVANCED"
}

func (s *CrossoverStrategyAdvanced) Evaluate(candles []entity.Candle, ctx *entity.StrategyContext) string {
	if len(candles) < 26 {
		return "HOLD"
	}

	prices := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
	}

	ma9 := indicators.MovingAverage(prices, 9)
	ma26 := indicators.MovingAverage(prices, 26)
	rsi := indicators.RSI(prices, 14)
	currentPrice := prices[len(prices)-1]

	if ctx.PositionQuantity == 0 {
		if ma9 > ma26 && currentPrice > ma9 && rsi < 70 {
			return "BUY"
		}
	} else {
		emaTrailing := indicators.MovingAverage(prices, 5)
		rsiPrev := indicators.RSI(prices[:len(prices)-1], 14)
		atr := indicators.ATRFromCandles(candles)
		stopLossThreshold := ctx.LastEntryPrice + atr*1.5

		stopLossHit := currentPrice < stopLossThreshold
		priceBelowTrailing := currentPrice < emaTrailing
		rsiReversal := rsiPrev > 80 && rsi < rsiPrev

		if stopLossHit || priceBelowTrailing || rsiReversal || (ma9 < ma26 && currentPrice < ma9) {
			return "SELL"
		}
	}

	return "HOLD"
}
