// internal/server/ws/ws_event.go

package ws

type Event struct {
	Type   string      `json:"type"`   // "candle" ou "decision"
	Symbol string      `json:"symbol"` // Ex: "BTCUSDT"
	Data   interface{} `json:"data"`   // Conteúdo do evento
}
