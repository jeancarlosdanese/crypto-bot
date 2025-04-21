// internal/domain/entity/indicator_snapshot.go

package entity

type IndicatorSnapshot struct {
	Timestamp  int64
	Price      float64
	Volume     float64
	EMAs       map[int]float64
	MACD       float64
	MACDSignal float64
	RSI        float64
	ATR        float64
	Volatility float64
	BBUpper    float64
	BBLower    float64
	BBWidth    float64
	Meta       map[string]any
}
