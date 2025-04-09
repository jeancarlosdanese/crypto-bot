// internal/domain/entity/position.go

package entity

type OpenPosition struct {
	Symbol     string       `bson:"symbol"`
	Interval   string       `bson:"interval"`
	EntryPrice float64      `bson:"entry_price"`
	Timestamp  int64        `bson:"timestamp"`
	Strategy   StrategyInfo `bson:"strategy"`
}
