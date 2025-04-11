// internal/domain/entity/decision_log.go

package entity

import (
	"time"

	"github.com/google/uuid"
)

type DecisionLog struct {
	BotID      uuid.UUID          `json:"bot_id"`
	Symbol     string             `json:"symbol"`
	Interval   string             `json:"interval"`
	Timestamp  int64              `json:"timestamp"`
	Decision   string             `json:"decision"`
	Indicators map[string]float64 `json:"indicators"`
	Context    map[string]any     `json:"context"`
	Strategy   StrategyInfo       `json:"strategy"`
	CreatedAt  time.Time          `json:"created_at"`
}
