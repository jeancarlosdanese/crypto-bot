// internal/domain/strategy/rsi2.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type RSI2Strategy struct{}

func (s *RSI2Strategy) Name() string {
	return "RSI2"
}

func (s *RSI2Strategy) Evaluate(candles []entity.Candle, ctx *entity.StrategyContext) string {
	if len(candles) < 3 {
		return "HOLD"
	}

	prices := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
	}

	rsi := indicators.RSI(prices, 2)

	if ctx.PositionQuantity == 0 && rsi < 10 {
		return "BUY"
	}

	if ctx.PositionQuantity > 0 && rsi > 90 {
		return "SELL"
	}

	return "HOLD"
}
