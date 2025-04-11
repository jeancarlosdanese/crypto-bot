package entity

import "github.com/google/uuid"

type Bot struct {
	ID           uuid.UUID `json:"id"`
	AccountID    uuid.UUID `json:"account_id"`
	Symbol       string    `json:"symbol"`
	Interval     string    `json:"interval"`
	StrategyName string    `json:"strategy_name"`
	Autonomous   bool      `json:"autonomous"`
	Active       bool      `json:"active"`
}
