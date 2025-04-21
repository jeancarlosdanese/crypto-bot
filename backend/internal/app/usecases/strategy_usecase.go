// internal/app/usecases/strategy_usecase.go

package usecases

import (
	"fmt"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/config"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/strategy"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	reporter "github.com/jeancarlosdanese/crypto-bot/internal/report"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"

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
	IndicatorService    *services.IndicatorService        // Serviço de indicadores para cálculos adicionais
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
	indicatorService *services.IndicatorService,
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
		IndicatorService: indicatorService,
		WindowSize:       windowSize,
		CandlesWindow:    make([]entity.Candle, 0, windowSize),
		LastDecision:     "HOLD",
	}
}

// ExecuteDecision avalia os indicadores já calculados e toma uma decisão (BUY, SELL ou HOLD).
func (s *StrategyUseCase) ExecuteDecision(timestamp int64) string {
	// Carrega configurações personalizadas do bot
	cfg, err := config.UnmarshalBotIndicatorConfig(s.Bot.ConfigJSON)
	if err != nil {
		logger.Error("❌ Erro ao carregar config_json", err)
		return "HOLD"
	}

	// Gera snapshot de indicadores técnicos com base nos candles e configurações
	snapshot := s.IndicatorService.GenerateSnapshot(s.CandlesWindow, cfg)
	if snapshot == nil {
		return "HOLD"
	}

	currentPrice := snapshot.Price
	rsi := snapshot.RSI
	atr := snapshot.ATR
	volatility := snapshot.Volatility

	// Contexto para a strategy
	strategyCtx := &entity.StrategyContext{
		Account:          s.Account,
		Bot:              s.Bot.Bot,
		PositionQuantity: s.PositionQuantity,
		LastEntryPrice:   s.LastEntryPrice,
		LastEntryTime:    s.LastEntryTimestamp,
	}

	// Estratégia toma a decisão com base no snapshot e contexto
	decision := s.Strategy.Evaluate(snapshot, strategyCtx)

	strategyName := s.Strategy.Name()
	strategyVersion := "1.0.1"

	// Filtros e execução
	switch decision {
	case "BUY":
		// Filtros técnicos mínimos
		if volatility < 0.1 {
			logger.Debug("🚫 Entrada bloqueada por baixa volatilidade", "volatility", volatility)
			return "HOLD"
		}
		if atr < 0.01 {
			logger.Debug("🚫 Entrada bloqueada por ATR insuficiente", "atr", atr)
			return "HOLD"
		}

		// Execução da entrada
		s.PositionQuantity = 1
		s.LastEntryPrice = currentPrice
		s.LastEntryTimestamp = timestamp
		s.LastDecision = "BUY"

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Save(entity.OpenPosition{
				BotID:      s.Bot.ID,
				EntryPrice: currentPrice,
				Timestamp:  timestamp,
			})
		}

		s.saveDecisionLog(strategyName, strategyVersion, "BUY", timestamp, nil, nil, nil)

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
		emaTrailing := snapshot.EMAs[cfg.GetTrailingEMA()]
		rsiPrev := indicators.RSIFromSnapshot(s.CandlesWindow, cfg.RSIPeriod, -1) // opcional: RSI anterior
		atrMultiplier := 1.5
		stopLossThreshold := s.LastEntryPrice + atr*atrMultiplier

		reason := "Desconhecido"
		if currentPrice < stopLossThreshold {
			reason = fmt.Sprintf("ATR stop hit (%.2f < %.2f)", currentPrice, stopLossThreshold)
		} else if currentPrice < emaTrailing {
			reason = fmt.Sprintf("Price < EMA (%d)", cfg.GetTrailingEMA())
		} else if rsiPrev > cfg.RSISell && rsi < rsiPrev {
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

		s.saveDecisionLog(strategyName, strategyVersion, "SELL", timestamp, nil, nil, nil)

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
