// internal/domain/entity/execution_log.go

package entity

import "time"

type ExecutionLog struct {
	Symbol    string       `bson:"symbol"`
	Interval  string       `bson:"interval"`
	Entry     TradePoint   `bson:"entry"`
	Exit      TradePoint   `bson:"exit"`
	Duration  int64        `bson:"duration_seconds"` // segundos entre entrada e saída
	Profit    float64      `bson:"profit"`           // lucro absoluto
	ROIPct    float64      `bson:"roi_pct"`          // retorno percentual
	Strategy  StrategyInfo `bson:"strategy"`
	CreatedAt time.Time    `bson:"created_at"`
}

type TradePoint struct {
	Price     float64 `bson:"price"`
	Timestamp int64   `bson:"timestamp"`
}
