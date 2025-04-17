// internal/factory/strategy_factory.go

package factory

import (
	"fmt"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/strategy"
)

func NewStrategyByName(name string) (strategy.Strategy, error) {
	switch name {
	case "CROSSOVER":
		return &strategy.CrossoverStrategy{}, nil
	case "CROSSOVER_ADVANCED":
		return &strategy.CrossoverStrategyAdvanced{}, nil
	case "EMA_FAN":
		return &strategy.EMAFanStrategy{}, nil
	default:
		return nil, fmt.Errorf("estratégia desconhecida: %s", name)
	}
}
