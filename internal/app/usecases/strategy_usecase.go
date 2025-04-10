// internal/app/usecases/strategy_usecase.go

package usecases

import (
	"fmt"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	reporter "github.com/jeancarlosdanese/crypto-bot/internal/report"
	service "github.com/jeancarlosdanese/crypto-bot/internal/services"
)

type StrategyUseCase struct {
	exchangeService       service.ExchangeService           // Serviço de exchange para obter dados de mercado
	decisionLogRepo       repository.DecisionLogRepository  // Repositório para registrar decisões
	executionLogRepo      repository.ExecutionLogRepository // Repositório para registrar execuções
	positionRepo          repository.PositionRepository     // Repositório para gerenciar posições abertas
	CandlesWindow         []entity.Candle                   // Janela de candles para análise
	WindowSize            int                               // Tamanho da janela de candles
	PositionQuantity      float64                           // Quantidade de posição atual (0 significa que não há posição)
	LastEntryPrice        float64                           // Último preço de entrada
	LastEntryTimestamp    int64                             // Último timestamp de entrada
	LastDecision          string                            // Última decisão tomada (BUY, SELL ou HOLD)
	TotalCandles          int                               // Contador global de candles processados
	LastCalibrationGlobal int                               // Valor global de TotalCandles no momento da calibração
}

// NewStrategyUseCase cria uma nova instância do StrategyUseCase com o tamanho de janela desejado.
func NewStrategyUseCase(exchangeService service.ExchangeService, decisionLogRepo repository.DecisionLogRepository, executionLogRepo repository.ExecutionLogRepository, positionRepo repository.PositionRepository, windowSize int) *StrategyUseCase {
	return &StrategyUseCase{
		exchangeService:  exchangeService,
		decisionLogRepo:  decisionLogRepo,
		executionLogRepo: executionLogRepo,
		positionRepo:     positionRepo,
		CandlesWindow:    make([]entity.Candle, 0, windowSize),
		WindowSize:       windowSize,
		LastDecision:     "HOLD",
	}
}

// UpdateCandle atualiza a janela de candles com o novo candle recebido.
func (d *StrategyUseCase) UpdateCandle(candle entity.Candle) {
	d.CandlesWindow = append(d.CandlesWindow, candle)
	d.TotalCandles++
	if len(d.CandlesWindow) > d.WindowSize {
		// Remove o candle mais antigo se a janela estiver cheia
		d.CandlesWindow = d.CandlesWindow[1:]
	}
}

// closingPrices extrai os preços de fechamento da janela de candles.
func (d *StrategyUseCase) closingPrices() []float64 {
	prices := make([]float64, len(d.CandlesWindow))
	for i, c := range d.CandlesWindow {
		prices[i] = c.Close
	}
	return prices
}

// EvaluateCrossover avalia o cruzamento de médias móveis e outros indicadores para tomar decisões de compra ou venda.
func (d *StrategyUseCase) EvaluateCrossover(symbol, interval string, timestamp int64) string {
	prices := d.closingPrices()
	if len(prices) < 26 {
		return "HOLD"
	}

	candlesSinceCalibration := d.TotalCandles - d.LastCalibrationGlobal
	if candlesSinceCalibration < 26 {
		logger.Debug("Aguardando formação de pelo menos 26 candles desde a calibração", "candles_formados", candlesSinceCalibration)
		return "HOLD"
	}

	ma9 := indicators.MovingAverage(prices, 9)
	ma26 := indicators.MovingAverage(prices, 26)
	volatility := indicators.Volatility(prices)
	rsi := indicators.RSI(prices, 14)
	atr := indicators.ATRFromCandles(d.CandlesWindow)
	currentPrice := prices[len(prices)-1]

	basicSignal := "HOLD"
	if ma9 > ma26 && currentPrice > ma9 && rsi < 70 {
		basicSignal = "BUY"
	} else if ma9 < ma26 && currentPrice < ma9 {
		basicSignal = "SELL"
	}

	indicatorsMap := map[string]float64{
		"ma9":        ma9,
		"ma26":       ma26,
		"volatility": volatility,
		"rsi":        rsi,
		"atr":        atr,
		"price":      currentPrice,
	}

	parameters := map[string]any{
		"ma_short":      9,
		"ma_long":       26,
		"rsi_threshold": 70,
		"atr_min":       0.5,
	}

	context := map[string]any{
		"window_size":      d.WindowSize,
		"total_candles":    d.TotalCandles,
		"last_calibration": d.LastCalibrationGlobal,
	}

	if d.PositionQuantity == 0 && basicSignal == "BUY" {
		d.PositionQuantity = 1
		d.LastEntryPrice = currentPrice
		d.LastEntryTimestamp = timestamp
		d.LastDecision = "BUY"

		logger.Info("Posição comprada", "preco", currentPrice, "timestamp", timestamp)

		// salva posição em aberto
		if d.positionRepo != nil {
			_ = d.positionRepo.Save(entity.OpenPosition{
				Symbol:     symbol,
				Interval:   interval,
				EntryPrice: currentPrice,
				Timestamp:  timestamp,
				Strategy: entity.StrategyInfo{
					Name:    "EvaluateCrossover",
					Version: "1.0.0",
				},
			})
		}

		d.saveDecisionLog("EvaluateCrossover", "1.0.0", symbol, interval, timestamp, "BUY", indicatorsMap, parameters, context)
		return "BUY"

	} else if d.PositionQuantity > 0 {
		// Lógica de saída inteligente
		ema5 := indicators.MovingAverage(prices, 5)
		rsiNow := indicators.RSI(prices, 14)
		rsiPrev := indicators.RSI(prices[:len(prices)-1], 14)
		atr := indicators.ATRFromCandles(d.CandlesWindow)
		stopLossThreshold := d.LastEntryPrice + atr*1.5
		priceBelowEma5 := currentPrice < ema5
		rsiReversal := rsiPrev > 80 && rsiNow < rsiPrev
		stopLossHit := currentPrice < stopLossThreshold

		if priceBelowEma5 || rsiReversal || stopLossHit || basicSignal == "SELL" {
			d.PositionQuantity = 0
			d.LastDecision = "SELL"

			logger.Info("💡 Critério de saída atingido", "preco", currentPrice, "timestamp", timestamp)

			if d.positionRepo != nil {
				_ = d.positionRepo.Delete(symbol)
			}

			d.saveDecisionLog("EvaluateCrossover", "1.0.0", symbol, interval, timestamp, "SELL", indicatorsMap, parameters, context)

			profit := currentPrice - d.LastEntryPrice
			duration := (timestamp - d.LastEntryTimestamp) / 1000
			roi := (profit / d.LastEntryPrice) * 100

			execLog := entity.ExecutionLog{
				Symbol:   symbol,
				Interval: interval,
				Entry: entity.TradePoint{
					Price:     d.LastEntryPrice,
					Timestamp: d.LastEntryTimestamp,
				},
				Exit: entity.TradePoint{
					Price:     currentPrice,
					Timestamp: timestamp,
				},
				Duration: duration,
				Profit:   profit,
				ROIPct:   roi,
				Strategy: entity.StrategyInfo{
					Name:    "EvaluateCrossover",
					Version: "1.0.1", // versão incrementada
				},
			}

			if d.executionLogRepo != nil {
				_ = d.executionLogRepo.Save(execLog)
			}

			go reporter.PrintExecutionSummary(d.executionLogRepo)

			logger.Info("💰 Execução registrada", "profit", profit, "roi_pct", roi, "duration", duration)
			return "SELL"
		}
	}

	return "HOLD"
}

// CalibrateLastEntry percorre a janela de candles para encontrar o último ponto de reversão significativo.
func (d *StrategyUseCase) CalibrateLastEntry() {
	prices := d.closingPrices()
	var lastSignal string = "HOLD"
	// Percorre os preços de trás para frente
	for i := len(prices) - 1; i > 0; i-- {
		window := prices[i:]
		signal := d.calculateSignal(window)
		if lastSignal != "HOLD" && signal != lastSignal && signal != "HOLD" {
			d.LastEntryPrice = prices[i]
			d.LastDecision = signal
			// Registra o TotalCandles atual como ponto de calibração
			d.LastCalibrationGlobal = d.TotalCandles - (len(prices) - i)
			logger.Debug("Ponto de reversão calibrado", "Sinal", signal, "Preço", prices[i], "Calibração Global", d.LastCalibrationGlobal)
			return
		}
		lastSignal = signal
	}
	d.LastEntryPrice = prices[len(prices)-1]
	d.LastDecision = "HOLD"
	d.LastCalibrationGlobal = d.TotalCandles
	logger.Debug("Nenhum ponto de reversão significativo encontrado, definindo estado default", "Calibração Global", d.LastCalibrationGlobal)
}

func (d *StrategyUseCase) EvaluateEMAFanWithVolume(symbol, interval string, timestamp int64) string {
	prices := d.closingPrices()
	if len(prices) < 40 || len(d.CandlesWindow) < 2 {
		return "HOLD"
	}

	// candlesSinceCalibration := d.TotalCandles - d.LastCalibrationGlobal
	// if candlesSinceCalibration < 40 {
	// 	logger.Debug("Aguardando formação de pelo menos 40 candles desde a calibração", "candles_formados", candlesSinceCalibration)
	// 	return "HOLD"
	// }

	// Calcula as EMAs
	emas := []float64{}
	periods := []int{10, 15, 20, 25, 30, 35, 40}
	for _, p := range periods {
		emas = append(emas, indicators.MovingAverage(prices, p))
	}

	// Verifica se todas as EMAs estão alinhadas de forma crescente
	isAligned := true
	for i := 1; i < len(emas); i++ {
		if emas[i] <= emas[i-1] {
			isAligned = false
			break
		}
	}

	slopeMap := calculateEMASlopes(prices, periods)

	for _, slope := range slopeMap {
		if slope < 0.08 { // mínimo de 0.05% de inclinação
			return "HOLD"
		}
	}

	// Verifica volume crescente comparando último volume com a média dos 10 anteriores
	lastVolume := d.CandlesWindow[len(d.CandlesWindow)-1].Volume
	avgVolume := 0.0
	for i := len(d.CandlesWindow) - 11; i < len(d.CandlesWindow)-1; i++ {
		avgVolume += d.CandlesWindow[i].Volume
	}
	avgVolume /= 10
	volumeConfirmed := lastVolume > avgVolume

	currentPrice := prices[len(prices)-1]

	indicatorsMap := map[string]float64{
		"price":        currentPrice,
		"last_volume":  lastVolume,
		"avg_volume":   avgVolume,
		"volume_ratio": lastVolume / avgVolume,
		"ema10":        emas[0],
		"ema40":        emas[len(emas)-1],
	}
	for i, p := range periods {
		indicatorsMap[fmt.Sprintf("ema%d", p)] = emas[i]
	}

	parameters := map[string]any{
		"emas":             periods,
		"volume_period":    10,
		"volume_ratio_min": 1.0,
	}

	context := map[string]any{
		"window_size":      d.WindowSize,
		"total_candles":    d.TotalCandles,
		"last_calibration": d.LastCalibrationGlobal,
	}

	strategyName := "EvaluateEMAFanWithVolume"
	strategyVersion := "1.0.0"

	// Estratégia de entrada
	if isAligned && volumeConfirmed && d.PositionQuantity == 0 {
		d.PositionQuantity = 1
		d.LastEntryPrice = currentPrice
		d.LastEntryTimestamp = timestamp
		d.LastDecision = "BUY"

		logger.Info("Posição comprada (EMA Fan)", "preco", currentPrice, "timestamp", timestamp)

		if d.positionRepo != nil {
			_ = d.positionRepo.Save(entity.OpenPosition{
				Symbol:     symbol,
				Interval:   interval,
				EntryPrice: currentPrice,
				Timestamp:  timestamp,
				Strategy: entity.StrategyInfo{
					Name:    strategyName,
					Version: strategyVersion,
				},
			})
		}

		d.saveDecisionLog(strategyName, strategyVersion, symbol, interval, timestamp, "BUY", indicatorsMap, parameters, context)
		return "BUY"
	}

	// Estratégia de saída
	if d.PositionQuantity > 0 && !isAligned {
		d.PositionQuantity = 0
		d.LastDecision = "SELL"

		logger.Info("Posição vendida (EMA Fan)", "preco", currentPrice, "timestamp", timestamp)

		if d.positionRepo != nil {
			_ = d.positionRepo.Delete(symbol)
		}

		d.saveDecisionLog(strategyName, strategyVersion, symbol, interval, timestamp, "SELL", indicatorsMap, parameters, context)

		// Salva execução
		profit := currentPrice - d.LastEntryPrice
		duration := (timestamp - d.LastEntryTimestamp) / 1000
		roi := (profit / d.LastEntryPrice) * 100

		execLog := entity.ExecutionLog{
			Symbol:   symbol,
			Interval: interval,
			Entry: entity.TradePoint{
				Price:     d.LastEntryPrice,
				Timestamp: d.LastEntryTimestamp,
			},
			Exit: entity.TradePoint{
				Price:     currentPrice,
				Timestamp: timestamp,
			},
			Duration: duration,
			Profit:   profit,
			ROIPct:   roi,
			Strategy: entity.StrategyInfo{
				Name:    strategyName,
				Version: strategyVersion,
			},
		}
		if d.executionLogRepo != nil {
			_ = d.executionLogRepo.Save(execLog)
		}

		go reporter.PrintExecutionSummary(d.executionLogRepo)
		return "SELL"
	}

	return "HOLD"
}

// calculateEMASlopes calcula a inclinação percentual das EMAs nos últimos dois pontos.
func calculateEMASlopes(prices []float64, periods []int) map[int]float64 {
	slopeMap := make(map[int]float64)
	for _, p := range periods {
		if len(prices) < p+2 {
			slopeMap[p] = 0
			continue
		}
		emaNow := indicators.MovingAverage(prices[len(prices)-1:], p)
		emaPrev := indicators.MovingAverage(prices[len(prices)-2:], p)
		if emaPrev != 0 {
			slopeMap[p] = ((emaNow - emaPrev) / emaPrev) * 100 // porcentagem
		} else {
			slopeMap[p] = 0
		}
	}
	return slopeMap
}

// calculateSignal aplica uma lógica simples baseada em MA9 e MA26 para determinar o sinal.
func (d *StrategyUseCase) calculateSignal(window []float64) string {
	if len(window) < 26 {
		return "HOLD"
	}
	ma9 := indicators.MovingAverage(window, 9)
	ma26 := indicators.MovingAverage(window, 26)
	currentPrice := window[len(window)-1]

	if ma9 > ma26 && currentPrice > ma9 {
		return "BUY"
	} else if ma9 < ma26 && currentPrice < ma9 {
		return "SELL"
	}
	return "HOLD"
}

func (d *StrategyUseCase) saveDecisionLog(strategyName, strategyVersin, symbol, interval string, timestamp int64, decision string, indicators map[string]float64, parameters, context map[string]any) {
	if d.decisionLogRepo == nil {
		return
	}

	log := entity.DecisionLog{
		Symbol:        symbol,
		Interval:      interval,
		Timestamp:     timestamp,
		Decision:      decision,
		PositionOpen:  d.PositionQuantity > 0,
		CandlesWindow: d.CandlesWindow,
		Indicators:    indicators,
		Strategy: entity.StrategyInfo{
			Name:       strategyName,
			Version:    strategyVersin,
			Parameters: parameters,
		},
		Context:   context,
		CreatedAt: time.Now(),
	}

	_ = d.decisionLogRepo.Save(log)
}
