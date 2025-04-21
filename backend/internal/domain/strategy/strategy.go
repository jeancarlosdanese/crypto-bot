// internal/domain/strategy/strategy.go

package strategy

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type Strategy interface {
	Name() string
	Evaluate(snapshot *entity.IndicatorSnapshot, ctx *entity.StrategyContext) string
}
