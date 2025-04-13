// internal/factory/exchange_factory.go

package factory

import (
	"fmt"
	"strings"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
	"github.com/jeancarlosdanese/crypto-bot/internal/services/binance"
)

type ExchangeFactory interface {
	NewExchangeService(exchangeName string, account *entity.Account) services.ExchangeService
}

type exchangeFactory struct{}

func NewExchangeFactory() ExchangeFactory {
	return &exchangeFactory{}
}

func (f *exchangeFactory) NewExchangeService(exchangeName string, account *entity.Account) services.ExchangeService {
	switch strings.ToLower(exchangeName) {
	case "binance":
		return binance.NewBinanceServiceWithKeys(*account.BinanceAPIKey, *account.BinanceAPISecret)
	default:
		panic(fmt.Sprintf("Exchange não suportada: %s", exchangeName))
	}
}
