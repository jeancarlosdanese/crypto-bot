// internal/app/usecases/strategy_usecase.go

package usecases

import (
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	service "github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type StrategyUseCase struct {
	Account             entity.Account                    // Conta do usuário
	Bot                 entity.Bot                        // Bot associado à conta
	Exchange            service.ExchangeService           // Serviço de exchange para obter dados de mercado
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
	bot entity.Bot,
	exchange service.ExchangeService,
	decisionRepo repository.DecisionLogRepository,
	executionRepo repository.ExecutionLogRepository,
	positionRepo repository.PositionRepository,
	windowSize int,
) *StrategyUseCase {
	return &StrategyUseCase{
		Account:          account,
		Bot:              bot,
		Exchange:         exchange,
		DecisionLogRepo:  decisionRepo,
		ExecutionLogRepo: executionRepo,
		PositionRepo:     positionRepo,
		WindowSize:       windowSize,
		CandlesWindow:    make([]entity.Candle, 0, windowSize),
		LastDecision:     "HOLD",
	}
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
