// internal/domain/entity/execution_log.go

package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExecutionLog struct {
	BotID     uuid.UUID    `json:"bot_id"`
	Symbol    string       `json:"symbol"`
	Interval  string       `json:"interval"`
	Entry     TradePoint   `json:"entry"`
	Exit      TradePoint   `json:"exit"`
	Duration  int64        `json:"duration"` // segundos entre entrada e saída
	Profit    float64      `json:"profit"`
	ROIPct    float64      `json:"roi_pct"`
	Strategy  StrategyInfo `json:"strategy"`
	CreatedAt time.Time    `json:"created_at"`
}

type TradePoint struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}
