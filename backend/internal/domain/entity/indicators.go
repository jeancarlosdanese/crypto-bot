// internal/domain/entity/indicator.go

package entity

type Indicators struct {
	Volatility float64 `bson:"volatility"`  // Volatilidade
	SMA        float64 `bson:"sma"`         // Média Móvel Simples
	EMA        float64 `bson:"ema"`         // Média Móvel Exponencial
	ATR        float64 `bson:"atr"`         // Average True Range
	RSI        float64 `bson:"rsi"`         // Índice de Força Relativa
	MacdLine   float64 `bson:"macd_line"`   // Linha MACD
	SignalLine float64 `bson:"signal_line"` // Linha de Sinal MACD
}
