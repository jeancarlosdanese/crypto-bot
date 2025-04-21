// internal/domain/strategy/crossover_advanced.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type CrossoverStrategyAdvanced struct{}

func (s *CrossoverStrategyAdvanced) Name() string {
	return "CROSSOVER_ADVANCED"
}

func (s *CrossoverStrategyAdvanced) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil {
		return "HOLD"
	}

	ema9 := snapshot.EMAs[9]
	ema26 := snapshot.EMAs[26]
	currentPrice := snapshot.Price
	rsi := snapshot.RSI
	atr := snapshot.ATR
	emaTrailing := snapshot.EMAs[5] // ou use o menor EMA com fallback

	if ctx.PositionQuantity == 0 {
		if ema9 > ema26 && currentPrice > ema9 && rsi < 70 {
			return "BUY"
		}
	} else {
		stopLossThreshold := ctx.LastEntryPrice + atr*1.5
		rsiReversal := rsi > 80 && rsi < ctx.LastEntryPrice // simplificação

		stopLossHit := currentPrice < stopLossThreshold
		priceBelowTrailing := currentPrice < emaTrailing

		if stopLossHit || priceBelowTrailing || rsiReversal || (ema9 < ema26 && currentPrice < ema9) {
			return "SELL"
		}
	}

	return "HOLD"
}

func (s *CrossoverStrategyAdvanced) EvaluateSnapshot(
	candles []entity.Candle,
	ctx *entity.StrategyContext,
	is *services.IndicatorService,
) string {
	snapshot := is.GenerateSnapshot(
		candles,
		[]int{9, 26},
		12, 26, 9,
		14,
		14,
		14,
		20,
	)
	return s.Evaluate(snapshot, ctx)
}
