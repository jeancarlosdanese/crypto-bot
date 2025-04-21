// internal/domain/config/bot_indicator_config.go

package config

import (
	"encoding/json"
	"fmt"
)

// BotIndicatorConfig representa os parâmetros técnicos configuráveis por bot.
type BotIndicatorConfig struct {
	EMAPeriods       []int      `json:"ema_periods"`
	MACD             MACDConfig `json:"macd"`
	RSIPeriod        int        `json:"rsi_period"`
	RSIBuy           float64    `json:"rsi_buy"`
	RSISell          float64    `json:"rsi_sell"`
	Bollinger        BBConfig   `json:"bollinger"`
	ATRPeriod        int        `json:"atr_period"`
	VolatilityWindow int        `json:"volatility_window"`
}

// MACDConfig define os períodos para o cálculo do MACD.
type MACDConfig struct {
	Short  int `json:"short"`
	Long   int `json:"long"`
	Signal int `json:"signal"`
}

// BBConfig define o período das bandas de Bollinger.
type BBConfig struct {
	Period int `json:"period"`
}

func (c *BotIndicatorConfig) GetTrailingEMA() int {
	if len(c.EMAPeriods) == 0 {
		return 9 // fallback padrão
	}
	min := c.EMAPeriods[0]
	for _, p := range c.EMAPeriods {
		if p < min {
			min = p
		}
	}
	return min
}

// UnmarshalBotIndicatorConfig converte o JSON do banco em um struct Go.
func UnmarshalBotIndicatorConfig(jsonData []byte) (*BotIndicatorConfig, error) {
	var cfg BotIndicatorConfig
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		return nil, fmt.Errorf("erro ao decodificar config_json: %w", err)
	}

	// Defaults de segurança (caso faltem campos)
	if len(cfg.EMAPeriods) == 0 {
		cfg.EMAPeriods = []int{9, 26}
	}
	if cfg.RSIPeriod == 0 {
		cfg.RSIPeriod = 14
	}
	if cfg.RSIBuy == 0 {
		cfg.RSIBuy = 10
	}
	if cfg.RSISell == 0 {
		cfg.RSISell = 90
	}
	if cfg.MACD.Short == 0 {
		cfg.MACD = MACDConfig{Short: 12, Long: 26, Signal: 9}
	}
	if cfg.Bollinger.Period == 0 {
		cfg.Bollinger.Period = 20
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = 14
	}
	if cfg.VolatilityWindow == 0 {
		cfg.VolatilityWindow = 14
	}

	return &cfg, nil
}
