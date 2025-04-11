package entity

import "github.com/google/uuid"

type OpenPosition struct {
	BotID      uuid.UUID `json:"bot_id"`
	EntryPrice float64   `json:"entry_price"`
	Timestamp  int64     `json:"timestamp"`
}