// internal/domain/strategy/strategy.go

package strategy

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type Strategy interface {
	Name() string
	Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string
	EvaluateSnapshot(candles []entity.Candle, ctx *entity.StrategyContext, is *services.IndicatorService) string
}
