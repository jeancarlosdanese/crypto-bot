package usecases

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
)

// CalibrateLastEntry recalibra o último ponto de entrada com base nos preços de fechamento.
func (s *StrategyUseCase) CalibrateLastEntry() {
	prices := s.ClosingPrices()
	var lastSignal string = "HOLD"
	for i := len(prices) - 1; i > 0; i-- {
		window := prices[i:]
		signal := s.calculateSignal(window)
		if lastSignal != "HOLD" && signal != lastSignal && signal != "HOLD" {
			s.LastEntryPrice = prices[i]
			s.LastDecision = signal
			s.LastCalibrationGlob = s.TotalCandles - (len(prices) - i)
			logger.Debug("Ponto de reversão calibrado", "Sinal", signal, "Preço", prices[i], "Calibração Global", s.LastCalibrationGlob)
			return
		}
		lastSignal = signal
	}
	s.LastEntryPrice = prices[len(prices)-1]
	s.LastDecision = "HOLD"
	s.LastCalibrationGlob = s.TotalCandles
	logger.Debug("Nenhum ponto de reversão significativo encontrado", "Calibração Global", s.LastCalibrationGlob)
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
