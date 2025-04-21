// internal/domain/strategy/crossover.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type CrossoverStrategy struct{}

func (s *CrossoverStrategy) Name() string {
	return "CROSSOVER"
}

func (s *CrossoverStrategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil {
		return "HOLD"
	}

	ema9 := snapshot.EMAs[9]
	ema26 := snapshot.EMAs[26]
	price := snapshot.Price
	rsi := snapshot.RSI

	if ctx.PositionQuantity == 0 {
		if ema9 > ema26 && price > ema9 && rsi < 70 {
			return "BUY"
		}
	} else {
		if ema9 < ema26 && price < ema9 {
			return "SELL"
		}
	}

	return "HOLD"
}

func (s *CrossoverStrategy) EvaluateSnapshot(
	candles []entity.Candle,
	ctx *entity.StrategyContext,
	is *services.IndicatorService,
) string {
	snapshot := is.GenerateSnapshot(
		candles,
		[]int{9, 26}, // EMA periods
		12, 26, 9,    // MACD config
		14, // RSI
		14, // ATR
		14, // Volatility
		20, // BB
	)
	return s.Evaluate(snapshot, ctx)
}
