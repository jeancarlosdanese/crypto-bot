// internal/services/binance/binance_service.go

package binance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	service "github.com/jeancarlosdanese/crypto-bot/internal/services"
)

var _ service.ExchangeService = (*BinanceService)(nil)

type BinanceService struct {
	client *binance.Client
}

// NewBinanceServiceWithKeys cria a instância a partir das chaves da conta.
func NewBinanceServiceWithKeys(apiKey, apiSecret string) *BinanceService {
	client := binance.NewClient(apiKey, apiSecret)
	return &BinanceService{client: client}
}

// GetCurrentPrice retorna o preço atual do símbolo.
func (s *BinanceService) GetCurrentPrice(symbol string) (float64, error) {
	price, err := s.client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil || len(price) == 0 {
		return 0, err
	}
	return strconv.ParseFloat(price[0].Price, 64)
}

// GetHistoricalCandles retorna um slice de candles para os últimos 'limit' candles.
func (s *BinanceService) GetHistoricalCandles(symbol string, interval string, limit int) ([]entity.Candle, error) {
	klines, err := s.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	var candles []entity.Candle
	for _, k := range klines {
		open, _ := strconv.ParseFloat(k.Open, 64)
		high, _ := strconv.ParseFloat(k.High, 64)
		low, _ := strconv.ParseFloat(k.Low, 64)
		closeVal, _ := strconv.ParseFloat(k.Close, 64)
		volume, _ := strconv.ParseFloat(k.Volume, 64)
		candles = append(candles, entity.Candle{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closeVal,
			Volume: volume,
			Time:   k.CloseTime / 1000,
		})
	}
	return candles, nil
}

// GetAccountPositions obtém e exibe os saldos da conta.
func (s *BinanceService) GetAccountPositions() error {
	account, err := s.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}
	// Exibe todos os ativos; você pode filtrar os relevantes (por exemplo, USDT)
	for _, asset := range account.Balances {
		// Exemplo de impressão para debug
		// Você pode usar log.Printf se preferir
		println("Ativo:", asset.Asset, "Disponível:", asset.Free, "Em Ordem:", asset.Locked)
	}
	return nil
}

func (s *BinanceService) GetBaseQuote(symbol string) (string, string, error) {
	info, err := s.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return "", "", err
	}

	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			return s.BaseAsset, s.QuoteAsset, nil
		}
	}
	return "", "", fmt.Errorf("símbolo não encontrado")
}

// PlaceBuyOrder envia uma ordem de compra limitada para o símbolo especificado.
func (s *BinanceService) PlaceBuyOrder(symbol string, quantity, price float64) (*binance.CreateOrderResponse, error) {
	return s.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(strconv.FormatFloat(quantity, 'f', -1, 64)).
		Price(strconv.FormatFloat(price, 'f', -1, 64)).
		Do(context.Background())
}

// PlaceSellOrder envia uma ordem de venda limitada para o símbolo especificado.
func (s *BinanceService) PlaceSellOrder(symbol string, quantity, price float64) (*binance.CreateOrderResponse, error) {
	return s.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(strconv.FormatFloat(quantity, 'f', -1, 64)).
		Price(strconv.FormatFloat(price, 'f', -1, 64)).
		Do(context.Background())
}
