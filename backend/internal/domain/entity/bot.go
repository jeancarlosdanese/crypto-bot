// internal/domain/entity/bot.go

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID         uuid.UUID `json:"id"`
	AccountID  uuid.UUID `json:"account_id"`
	StrategyID uuid.UUID `json:"strategy_id"`
	Symbol     string    `json:"symbol"`
	Interval   string    `json:"interval"`
	Autonomous bool      `json:"autonomous"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BotWithStrategy struct {
	Bot
	StrategyName string `json:"strategy_name"`
}
