// internal/app/usecases/strategy_usecase.go

package usecases

import (
	"time"

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
	// Contexto da estratégia (posição atual, preços, etc.)
	strategyCtx := &entity.StrategyContext{
		Account:          s.Account,
		Bot:              s.Bot.Bot,
		PositionQuantity: s.PositionQuantity,
		LastEntryPrice:   s.LastEntryPrice,
		LastEntryTime:    s.LastEntryTimestamp,
	}

	// A estratégia avalia e decide com base no snapshot gerado por ela
	decision := s.Strategy.EvaluateSnapshot(s.CandlesWindow, strategyCtx, s.IndicatorService)
	strategyName := s.Strategy.Name()
	strategyVersion := "1.0.1"

	if decision == "HOLD" {
		return "HOLD"
	}

	price := s.CandlesWindow[len(s.CandlesWindow)-1].Close

	switch decision {
	case "BUY":
		s.PositionQuantity = 1
		s.LastEntryPrice = price
		s.LastEntryTimestamp = timestamp
		s.LastDecision = "BUY"

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Save(entity.OpenPosition{
				BotID:      s.Bot.ID,
				EntryPrice: price,
				Timestamp:  timestamp,
			})
		}

		s.saveDecisionLog(strategyName, strategyVersion, "BUY", timestamp, nil, nil, nil)

		serverws.Publish(s.Bot.ID.String(), serverws.Event{
			Type: "decision",
			Data: map[string]interface{}{
				"time":     timestamp / 1000,
				"price":    price,
				"decision": "BUY",
			},
		})

		logger.Info("📈 Entrada executada", "symbol", s.Bot.Symbol, "price", price)
		return "BUY"

	case "SELL":
		s.PositionQuantity = 0
		s.LastDecision = "SELL"

		if s.PositionRepo != nil {
			_ = s.PositionRepo.Delete(s.Bot.ID)
		}

		profit := price - s.LastEntryPrice
		roi := (profit / s.LastEntryPrice) * 100
		duration := (timestamp - s.LastEntryTimestamp) / 1000

		s.saveDecisionLog(strategyName, strategyVersion, "SELL", timestamp, nil, nil, nil)

		if s.ExecutionLogRepo != nil {
			_ = s.ExecutionLogRepo.Save(entity.ExecutionLog{
				BotID:     s.Bot.ID,
				Symbol:    s.Bot.Symbol,
				Interval:  s.Bot.Interval,
				Entry:     entity.TradePoint{Price: s.LastEntryPrice, Timestamp: s.LastEntryTimestamp},
				Exit:      entity.TradePoint{Price: price, Timestamp: timestamp},
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
				"price":    price,
				"decision": "SELL",
			},
		})

		logger.Info("📉 Saída executada", "symbol", s.Bot.Symbol, "price", price, "roi", roi)
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
