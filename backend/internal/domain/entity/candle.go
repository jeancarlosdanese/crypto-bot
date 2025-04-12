// internal/domain/entity/candle.go

package entity

// Define o tipo Candle para armazenar os dados de um candle
type Candle struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Time   int64
}
