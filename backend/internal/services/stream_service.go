// internal/services/stream_service.go

package services

type StreamService interface {
	Start(symbol, interval string) error
	StartMany(pairs map[string]string) error
	Stop(symbol string)
	StopAll()
}
