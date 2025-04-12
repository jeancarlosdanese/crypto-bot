// internal/domain/dto/bot_dto.go

package dto

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type BotResponseDTO struct {
	ID           string `json:"id"`
	Symbol       string `json:"symbol"`
	Interval     string `json:"interval"`
	StrategyName string `json:"strategy_name"`
	Autonomous   bool   `json:"autonomous"`
	Active       bool   `json:"active"`
}

func NewBotResponseDTO(bot *entity.Bot) BotResponseDTO {
	return BotResponseDTO{
		ID:           bot.ID.String(),
		Symbol:       bot.Symbol,
		Interval:     bot.Interval,
		StrategyName: bot.StrategyName,
		Autonomous:   bot.Autonomous,
		Active:       bot.Active,
	}
}
