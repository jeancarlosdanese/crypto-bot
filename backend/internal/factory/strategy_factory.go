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
	case "RSI2":
		return &strategy.RSI2Strategy{}, nil

	// 🆕 Estratégias sugeridas para evolução:

	case "VOLUME_SPIKE":
		return &strategy.VolumeSpikeStrategy{}, nil

	case "BB_REBOUND":
		return &strategy.BollingerReboundStrategy{}, nil

	case "MACD_CROSS":
		return &strategy.MACDCrossStrategy{}, nil

	// case "AI":
	// 	return &strategy.AIStrategy{}, nil // Experimental (pode usar GPT ou modelo local)

	default:
		return nil, fmt.Errorf("estratégia desconhecida: %s", name)
	}
}
