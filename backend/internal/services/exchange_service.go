// internal/service/exchange_service.go

package services

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type ExchangeService interface {
	GetAccountPositions() error
	GetCurrentPrice(symbol string) (float64, error)
	GetHistoricalCandles(symbol string, interval string, limit int) ([]entity.Candle, error)
	GetBaseQuote(symbol string) (string, string, error)
}
