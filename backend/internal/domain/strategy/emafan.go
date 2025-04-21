// internal/domain/strategy/emafan.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type EMAFanStrategy struct{}

func (s *EMAFanStrategy) Name() string {
	return "EMA_FAN"
}

func (s *EMAFanStrategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil || len(snapshot.EMAs) < 7 {
		return "HOLD"
	}

	// Definição dos períodos do leque
	periods := []int{10, 15, 20, 25, 30, 35, 40}

	// Verifica alinhamento crescente das EMAs
	aligned := true
	prev := snapshot.EMAs[periods[0]]
	for _, p := range periods[1:] {
		curr := snapshot.EMAs[p]
		if curr <= prev {
			aligned = false
			break
		}
		prev = curr
	}

	if !aligned {
		if ctx.PositionQuantity > 0 {
			return "SELL"
		}
		return "HOLD"
	}

	// Confirmação por volume
	if snapshot.Volume > 0 && ctx.PositionQuantity == 0 {
		// No futuro: considerar média de volumes via snapshot também
		return "BUY"
	}

	return "HOLD"
}

func (s *EMAFanStrategy) EvaluateSnapshot(
	candles []entity.Candle,
	ctx *entity.StrategyContext,
	is *services.IndicatorService,
) string {
	snapshot := is.GenerateSnapshot(
		candles,
		[]int{10, 15, 20, 25, 30, 35, 40},
		0, 0, 0, // MACD não usado
		14,
		14,
		14,
		20,
	)
	return s.Evaluate(snapshot, ctx)
}
