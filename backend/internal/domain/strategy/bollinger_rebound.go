// internal/domain/strategy/bollinger_rebound.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type BollingerReboundStrategy struct{}

func (s *BollingerReboundStrategy) Name() string {
	return "BB_REBOUND"
}

func (s *BollingerReboundStrategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil {
		return "HOLD"
	}

	price := snapshot.Price
	lower := snapshot.BBLower
	upper := snapshot.BBUpper
	width := snapshot.BBWidth

	// Rejeita bandas muito estreitas (mercado parado)
	if width < 0.5 || lower == 0 || upper == 0 {
		return "HOLD"
	}

	if ctx.PositionQuantity == 0 && price < lower {
		return "BUY"
	}

	if ctx.PositionQuantity > 0 && price > upper {
		return "SELL"
	}

	return "HOLD"
}
