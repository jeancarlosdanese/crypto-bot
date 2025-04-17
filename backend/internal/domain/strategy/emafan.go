// internal/domain/strategy/emafan.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type EMAFanStrategy struct{}

func (s *EMAFanStrategy) Name() string {
	return "EMA_FAN"
}

func (s *EMAFanStrategy) Evaluate(candles []entity.Candle, ctx *entity.StrategyContext) string {
	if len(candles) < 40 {
		return "HOLD"
	}

	prices := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
	}

	periods := []int{10, 15, 20, 25, 30, 35, 40}
	aligned := true
	prev := indicators.MovingAverage(prices, periods[0])

	for _, p := range periods[1:] {
		curr := indicators.MovingAverage(prices, p)
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

	lastVolume := candles[len(candles)-1].Volume
	avgVolume := 0.0
	for i := len(candles) - 11; i < len(candles)-1; i++ {
		avgVolume += candles[i].Volume
	}
	avgVolume /= 10

	if lastVolume > avgVolume && ctx.PositionQuantity == 0 {
		return "BUY"
	}

	return "HOLD"
}
