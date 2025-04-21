// internal/domain/strategy/volume_spike.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type VolumeSpikeStrategy struct{}

func (s *VolumeSpikeStrategy) Name() string {
	return "VOLUME_SPIKE"
}

func (s *VolumeSpikeStrategy) Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string {
	if snapshot == nil || snapshot.Volume == 0 {
		return "HOLD"
	}

	// 🔍 Supondo que a média de volume já foi calculada no snapshot
	avgVolume, ok := snapshot.Meta["avg_volume"].(float64)
	if !ok || avgVolume == 0 {
		return "HOLD"
	}

	// 📈 Volume atual é maior que 2x a média dos últimos N candles
	volumeSpikeThreshold := 2.0

	if ctx.PositionQuantity == 0 && snapshot.Volume > avgVolume*volumeSpikeThreshold {
		return "BUY"
	}

	if ctx.PositionQuantity > 0 && snapshot.Volume < avgVolume {
		return "SELL"
	}

	return "HOLD"
}
