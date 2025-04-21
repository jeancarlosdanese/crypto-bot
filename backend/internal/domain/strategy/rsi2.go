// internal/domain/strategy/rsi2.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type RSI2Strategy struct{}

func (s *RSI2Strategy) Name() string {
	return "RSI2"
}

func (s *RSI2Strategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil {
		return "HOLD"
	}

	rsi := snapshot.RSI

	if ctx.PositionQuantity == 0 && rsi < 10 {
		return "BUY"
	}

	if ctx.PositionQuantity > 0 && rsi > 90 {
		return "SELL"
	}

	return "HOLD"
}
