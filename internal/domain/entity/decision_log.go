// internal/domain/entity/decision_log.go

package entity

import "time"

// DecisionLog representa uma decisão tomada pela estratégia
type DecisionLog struct {
	Symbol        string             `bson:"symbol"`
	Interval      string             `bson:"interval"`
	Timestamp     int64              `bson:"timestamp"` // timestamp do candle fechado
	Decision      string             `bson:"decision"`
	PositionOpen  bool               `bson:"position_active"`
	CandlesWindow []Candle           `bson:"candles_window"`
	Indicators    map[string]float64 `bson:"indicators"`
	Strategy      StrategyInfo       `bson:"strategy"`
	Context       map[string]any     `bson:"context"`
	CreatedAt     time.Time          `bson:"created_at"`
}
