// internal/domain/strategy/crossover.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type CrossoverStrategy struct{}

func (s *CrossoverStrategy) Name() string {
	return "CROSSOVER"
}

func (s *CrossoverStrategy) Evaluate(candles []entity.Candle, ctx *entity.StrategyContext) string {
	if len(candles) < 26 {
		return "HOLD"
	}

	prices := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
	}

	ma9 := indicators.MovingAverage(prices, 9)
	ma26 := indicators.MovingAverage(prices, 26)
	currentPrice := prices[len(prices)-1]

	if ctx.PositionQuantity == 0 {
		if ma9 > ma26 && currentPrice > ma9 {
			return "BUY"
		}
	} else {
		if ma9 < ma26 && currentPrice < ma9 {
			return "SELL"
		}
	}

	return "HOLD"
}
