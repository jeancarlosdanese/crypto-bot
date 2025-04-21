// internal/domain/strategy/macd_cross.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type MACDCrossStrategy struct{}

func (s *MACDCrossStrategy) Name() string {
	return "MACD_CROSS"
}

func (s *MACDCrossStrategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil {
		return "HOLD"
	}

	// Último valor do MACD e da linha de sinal
	macd := snapshot.MACD
	signal := snapshot.MACDSignal

	prevMACD, ok1 := snapshot.Meta["prev_macd"].(float64)
	prevSignal, ok2 := snapshot.Meta["prev_signal"].(float64)

	if !ok1 || !ok2 {
		return "HOLD"
	}

	// Cruzamento positivo: MACD cruza acima da linha de sinal
	if ctx.PositionQuantity == 0 && prevMACD < prevSignal && macd > signal {
		return "BUY"
	}

	// Cruzamento negativo: MACD cruza abaixo da linha de sinal
	if ctx.PositionQuantity > 0 && prevMACD > prevSignal && macd < signal {
		return "SELL"
	}

	return "HOLD"
}
