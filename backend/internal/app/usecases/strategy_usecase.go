// internal/app/usecases/strategy_usecase.go

package usecases

import (
	"fmt"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/strategy"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	reporter "github.com/jeancarlosdanese/crypto-bot/internal/report"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
	"github.com/jeancarlosdanese/crypto-bot/internal/utils"

	serverws "github.com/jeancarlosdanese/crypto-bot/internal/server/ws"
)

type StrategyUseCase struct {
	Account             entity.Account                    // Conta do usuário
	Bot                 entity.BotWithStrategy            // Bot associado à conta
	Exchange            services.ExchangeService          // Serviço de exchange para obter dados de mercado
	Strategy            strategy.Strategy                 // Implementação da estratégia de trading
	DecisionLogRepo     repository.DecisionLogRepository  // Repositório para registrar decisões
	ExecutionLogRepo    repository.ExecutionLogRepository // Repositório para registrar execuções
	PositionRepo        repository.PositionRepository     // Repositório para gerenciar posições abertas
	WindowSize          int                               // Tamanho da janela de candles
	CandlesWindow       []entity.Candle                   // Janela de candles para análise
	PositionQuantity    float64                           // Quantidade de posição atual (0 significa que não há posição)
	LastEntryPrice      float64                           // Último preço de entrada
	LastEntryTimestamp  int64                             // Último timestamp de entrada
	LastDecision        string                            // Última decisão tomada (BUY, SELL ou HOLD)
	TotalCandles        int                               // Contador global de candles processados
	LastCalibrationGlob int                               // Valor global de TotalCandles no momento da calibração
}

// NewStrategyUseCase cria uma nova instância do StrategyUseCase com o tamanho de janela desejado.
func NewStrategyUseCase(
	account entity.Account,
	bot entity.BotWithStrategy,
	exchange services.ExchangeService,
	strategyImpl strategy.Strategy,
	decisionRepo repository.DecisionLogRepository,
	executionRepo repository.ExecutionLogRepository,
	positionRepo repository.PositionRepository,
	windowSize int,
) *StrategyUseCase {
	return &StrategyUseCase{
		Account:          account,
		Bot:              bot,
		Exchange:         exchange,
		Strategy:         strategyImpl,
		DecisionLogRepo:  decisionRepo,
		ExecutionLogRepo: executionRepo,
		PositionRepo:     positionRepo,
		WindowSize:       windowSize,
		CandlesWindow:    make([]entity.Candle, 0, windowSize),
		LastDecision:     "HOLD",
	}
}

// ExecuteDecision avalia a janela de candles atual e toma uma decisão de trading (BUY, SELL ou HOLD).
func (s *StrategyUseCase) ExecuteDecision(timestamp int64) string {
	prices := s.ClosingPrices()
	if len(prices) < 26 {
		return "HOLD"
	}

	// 🔧 Indicadores calculados
	ma9 := indicators.MovingAverage(prices, 9)
	ma26 := indicators.MovingAverage(prices, 26)
	rsi := indicators.RSI(prices, 14)
	volatility := indicators.Volatility(prices)
	atr := indicators.ATRFromCandles(s.CandlesWindow)
	currentPrice := prices[len(prices)-1]

	// 🔧 Parâmetros configuráveis
	params := map[string]any{
		"ma_short":           9,
		"ma_long":            26,
		"rsi_threshold":      70,
		"volatility_min":     0.0,
		"atr_min":            0.0,
		"ema_trailing":       5,
		"rsi_exit_threshold": 80,
		"atr_multiplier":     1.5,
	}

	// 🔧 Contexto do bot
	ctx := map[string]any{
		"candles_total": s.TotalCandles,
		"calibrated_at": s.LastCalibrationGlob,
	}

	// 🔧 Indicadores
	indicatorsMap := map[string]float64{
		"ma9":        ma9,
		"ma26":       ma26,
		"rsi":        rsi,
		"volatility": volatility,
		"atr":        atr,
		"price":      currentPrice,
	}

	strategyVersion := "1.0.1"
	strategyName := s.Strategy.Name()

	// 🔍 Chamada à strategy
	strategyCtx := &entity.StrategyContext{
		Account:          s.Account,
		Bot:              s.Bot.Bot,
		PositionQuantity: s.PositionQuantity,
		LastEntryPrice:   s.LastEntryPrice,
		LastEntryTime:    s.LastEntryTimestamp,
	}
	decision := s.Strategy.Evaluate(s.CandlesWindow, strategyCtx)

	switch decision {
	case "BUY":
		// Filtros de entrada
		if volatility < utils.GetFloatParam(params, "volatility_min", 0.0) {
			logger.Debug("🚫 Entrada bloqueada por baixa volatilidade",
				"symbol", s.Bot.Symbol,
				"volatility", volatility,
			)
			return "HOLD"
		}
		if atr < utils.GetFloatParam(params, "atr_min", 0.0) {
			logger.Debug("🚫 Entrada bloqueada por ATR insuficiente",
				"symbol", s.Bot.Symbol,
				"atr", atr,
			)
			return "HOLD"
		}

		// Executa entrada
		s.PositionQuantity = 1
		s.LastEntryPrice = currentPrice
		s.LastEntryTimestamp = timestamp
		s.LastDecision = "BUY"

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Save(entity.OpenPosition{
				BotID:      s.Bot.ID,
				EntryPrice: s.LastEntryPrice,
				Timestamp:  timestamp,
			})
		}

		s.saveDecisionLog(strategyName, strategyVersion, "BUY", timestamp, indicatorsMap, params, ctx)

		serverws.Publish(s.Bot.ID.String(), serverws.Event{
			Type: "decision",
			Data: map[string]interface{}{
				"time":     timestamp / 1000,
				"price":    currentPrice,
				"decision": "BUY",
			},
		})

		logger.Info("📈 Entrada executada", "symbol", s.Bot.Symbol, "price", currentPrice)
		return "BUY"

	case "SELL":
		// Critérios técnicos de saída
		emaTrailing := indicators.MovingAverage(prices, utils.GetIntParam(params, "ema_trailing", 5))
		rsiPrev := indicators.RSI(prices[:len(prices)-1], 14)
		atrMultiplier := utils.GetFloatParam(params, "atr_multiplier", 1.5)
		stopLossThreshold := s.LastEntryPrice + atr*atrMultiplier

		reason := "Desconhecido"
		if currentPrice < stopLossThreshold {
			reason = fmt.Sprintf("ATR stop hit (%.2f < %.2f)", currentPrice, stopLossThreshold)
		} else if currentPrice < emaTrailing {
			reason = fmt.Sprintf("Price < EMA%d (%.2f < %.2f)", utils.GetIntParam(params, "ema_trailing", 5), currentPrice, emaTrailing)
		} else if rsiPrev > utils.GetFloatParam(params, "rsi_exit_threshold", 80) && rsi < rsiPrev {
			reason = fmt.Sprintf("RSI reversal (%.2f < %.2f)", rsi, rsiPrev)
		}

		s.PositionQuantity = 0
		s.LastDecision = "SELL"

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Delete(s.Bot.ID)
		}

		profit := currentPrice - s.LastEntryPrice
		roi := (profit / s.LastEntryPrice) * 100
		duration := (timestamp - s.LastEntryTimestamp) / 1000

		s.saveDecisionLog(strategyName, strategyVersion, "SELL", timestamp, indicatorsMap, params, ctx)

		if s.ExecutionLogRepo != nil {
			_ = s.ExecutionLogRepo.Save(entity.ExecutionLog{
				BotID:     s.Bot.ID,
				Symbol:    s.Bot.Symbol,
				Interval:  s.Bot.Interval,
				Entry:     entity.TradePoint{Price: s.LastEntryPrice, Timestamp: s.LastEntryTimestamp},
				Exit:      entity.TradePoint{Price: currentPrice, Timestamp: timestamp},
				Profit:    profit,
				ROIPct:    roi,
				Duration:  duration,
				Strategy:  entity.StrategyInfo{Name: strategyName, Version: strategyVersion},
				CreatedAt: time.Now(),
			})
		}

		serverws.Publish(s.Bot.ID.String(), serverws.Event{
			Type: "decision",
			Data: map[string]interface{}{
				"time":     timestamp / 1000,
				"price":    currentPrice,
				"decision": "SELL",
			},
		})

		logger.Info("📉 Saída executada", "symbol", s.Bot.Symbol, "price", currentPrice, "reason", reason)
		go reporter.PrintExecutionSummary(s.ExecutionLogRepo)
		return "SELL"
	}

	return "HOLD"
}

// UpdateCandle atualiza a janela de candles com o novo candle recebido.
func (s *StrategyUseCase) UpdateCandle(candle entity.Candle) {
	s.CandlesWindow = append(s.CandlesWindow, candle)
	s.TotalCandles++
	if len(s.CandlesWindow) > cap(s.CandlesWindow) {
		s.CandlesWindow = s.CandlesWindow[1:]
	}
}

// closingPrices extrai os preços de fechamento dos candles na janela atual.
func (s *StrategyUseCase) ClosingPrices() []float64 {
	prices := make([]float64, len(s.CandlesWindow))
	for i, c := range s.CandlesWindow {
		prices[i] = c.Close
	}
	return prices
}

func (s *StrategyUseCase) saveDecisionLog(strategy, version, decision string, timestamp int64, indicators map[string]float64, params, ctx map[string]any) {
	if s.DecisionLogRepo == nil {
		return
	}
	_ = s.DecisionLogRepo.Save(entity.DecisionLog{
		BotID:      s.Bot.ID,
		Symbol:     s.Bot.Symbol,
		Interval:   s.Bot.Interval,
		Timestamp:  timestamp,
		Decision:   decision,
		Indicators: indicators,
		Context:    ctx,
		Strategy: entity.StrategyInfo{
			Name:       strategy,
			Version:    version,
			Parameters: params,
		},
		CreatedAt: time.Now(),
	})
}
