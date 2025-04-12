// internal/services/binance/stream_service.go

package binance

import (
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	serverws "github.com/jeancarlosdanese/crypto-bot/internal/server/ws"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type BinanceStreamService interface {
	Start(symbol, interval string) error
	StartMany(pairs map[string]string) error
	Stop(symbol string)
	StopAll()
}

type binanceStreamService struct {
	strategy       *usecases.StrategyUseCase
	binanceService *BinanceService
	active         map[string]chan struct{}
}

func NewBinanceStreamService(strategy *usecases.StrategyUseCase, binanceService *BinanceService) BinanceStreamService {
	return &binanceStreamService{
		strategy:       strategy,
		binanceService: binanceService,
		active:         make(map[string]chan struct{}),
	}
}

func (b *binanceStreamService) Start(symbol, interval string) error {
	symbol = strings.ToLower(symbol)
	logger.Info("[StreamService] Iniciando monitoramento", "symbol", symbol, "interval", interval)

	candles, err := b.binanceService.GetHistoricalCandles(strings.ToUpper(symbol), interval, b.strategy.WindowSize)
	if err != nil {
		logger.Error("[StreamService] Erro ao obter candles históricos", err, "symbol", symbol)
		return err
	}

	for _, c := range candles {
		b.strategy.UpdateCandle(c)
	}

	b.strategy.CalibrateLastEntry()
	stopChan := make(chan struct{})
	b.active[symbol] = stopChan

	go func() {
		var current entity.Candle
		reconnectDelay := 5 * time.Second
		const maxUptime = 23*time.Hour + 55*time.Minute
		heartbeatTicker := time.NewTicker(5 * time.Minute)
		defer heartbeatTicker.Stop()

		for {
			timer := time.NewTimer(maxUptime)

			wsHandler := func(event *binance.WsKlineEvent) {
				k := event.Kline
				open, _ := strconv.ParseFloat(k.Open, 64)
				high, _ := strconv.ParseFloat(k.High, 64)
				low, _ := strconv.ParseFloat(k.Low, 64)
				closeVal, _ := strconv.ParseFloat(k.Close, 64)
				volume, _ := strconv.ParseFloat(k.Volume, 64)

				if current.Open == 0 {
					current = entity.Candle{
						Open:   open,
						High:   high,
						Low:    low,
						Close:  closeVal,
						Volume: volume,
					}
				} else {
					if high > current.High {
						current.High = high
					}
					if low < current.Low {
						current.Low = low
					}
					current.Close = closeVal
				}

				if k.IsFinal {
					b.strategy.UpdateCandle(current)

					// 🔥 Publicar candle no WebSocket
					serverws.Publish(b.strategy.Bot.ID.String(), serverws.Event{
						Type: "candle",
						Data: map[string]interface{}{
							"time":  k.EndTime / 1000, // frontend espera timestamp em segundos
							"open":  current.Open,
							"high":  current.High,
							"low":   current.Low,
							"close": current.Close,
						},
					})

					// timestamp do candle finalizado (já vem como int64 da Binance)
					decision := b.strategy.EvaluateCrossover(k.EndTime)

					if decision != "HOLD" {
						logger.Info("[StreamService] Decisão tomada",
							"symbol", symbol,
							"interval", interval,
							"decision", decision,
						)
					}

					current = entity.Candle{}
				}
			}

			errHandler := func(err error) {
				logger.Error("[StreamService] Erro no WebSocket", err, "symbol", symbol)
			}

			done, _, err := binance.WsKlineServe(symbol, interval, wsHandler, errHandler)
			if err != nil {
				logger.Warn("[StreamService] Erro ao conectar. Tentando reconectar...", "symbol", symbol, "espera", reconnectDelay)
				time.Sleep(reconnectDelay)
				continue
			}

			startTime := time.Now()

			run := true
			for run {
				select {
				case <-done:
					logger.Warn("[StreamService] Conexão encerrada. Reconnectando...", "symbol", symbol)
					run = false
				case <-timer.C:
					logger.Warn("[StreamService] Reconectando após tempo máximo de conexão", "symbol", symbol)
					run = false
				case <-heartbeatTicker.C:
					uptime := time.Since(startTime).Round(time.Second)
					logger.Debug("[StreamService] Conexão viva", "symbol", symbol, "uptime", uptime.String())
				case <-stopChan:
					logger.Info("[StreamService] Stream parada manualmente", "symbol", symbol)
					timer.Stop()
					return
				}
			}

			timer.Stop()
			time.Sleep(reconnectDelay)
		}
	}()

	return nil
}

// StartMany inicia múltiplas streams de forma assíncrona
// pairs é um mapa onde a chave é o símbolo e o valor é o intervalo
// Exemplo:
//
//	streamService.StartMany(map[string]string{
//	    "btcusdt": "1m",
//	    "ethusdt": "5m",
//	    "bnbusdt": "1m",
//	})
func (b *binanceStreamService) StartMany(pairs map[string]string) error {
	for symbol, interval := range pairs {
		err := b.Start(symbol, interval)
		if err != nil {
			logger.Error("[StreamService] Erro ao iniciar stream", err, "symbol", symbol, "interval", interval)
		}
	}
	return nil
}

func (b *binanceStreamService) Stop(symbol string) {
	if ch, ok := b.active[symbol]; ok {
		close(ch)
		delete(b.active, symbol)
		logger.Info("[StreamService] Stream parado", "symbol", symbol)
	}
}

func (b *binanceStreamService) StopAll() {
	for symbol := range b.active {
		b.Stop(symbol)
	}
}

var _ services.StreamService = (*binanceStreamService)(nil)
