// internal/services/indicator_service.go

package services

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/config"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type IndicatorService struct{}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{}
}

func (s *IndicatorService) GenerateSnapshot(
	candles []entity.Candle,
	cfg *config.BotIndicatorConfig,
) *entity.IndicatorSnapshot {
	if len(candles) < 2 {
		return nil
	}

	prices := make([]float64, len(candles))
	volumes := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
		volumes[i] = c.Volume
	}

	// ✅ EMAs dinâmicas
	emaMap := make(map[int]float64)
	for _, period := range cfg.EMAPeriods {
		emaMap[period] = indicators.LastEMA(prices, period)
	}

	// ✅ MACD (últimos valores)
	macd, signal, _ := indicators.MACD(prices, cfg.MACD.Short, cfg.MACD.Long, cfg.MACD.Signal)
	var macdVal, macdSignal float64
	if len(macd) > 0 {
		macdVal = macd[len(macd)-1]
	}
	if len(signal) > 0 {
		macdSignal = signal[len(signal)-1]
	}

	// ✅ RSI
	rsi := indicators.RSI(prices, cfg.RSIPeriod)

	// ✅ ATR
	atr := indicators.ATRFromCandles(candles[len(candles)-cfg.ATRPeriod:])

	// ✅ Volatilidade
	volatility := indicators.Volatility(prices[len(prices)-cfg.VolatilityWindow:])

	// ✅ Bollinger Bands
	bbUpper, bbLower := indicators.BollingerBands(prices, cfg.Bollinger.Period)
	bbWidth := bbUpper - bbLower

	// ✅ Meta: Preço atual em relação à média de volume dos últimos 10 candles
	avgVolume := indicators.AverageVolume(candles[len(candles)-11:]) // últimos 10 candles

	// ✅ Snapshot final
	return &entity.IndicatorSnapshot{
		Timestamp:  candles[len(candles)-1].Time,
		Price:      prices[len(prices)-1],
		Volume:     volumes[len(volumes)-1],
		EMAs:       emaMap,
		MACD:       macdVal,
		MACDSignal: macdSignal,
		RSI:        rsi,
		ATR:        atr,
		Volatility: volatility,
		BBUpper:    bbUpper,
		BBLower:    bbLower,
		BBWidth:    bbWidth,
		Meta: map[string]any{
			"avg_volume":  avgVolume,
			"prev_macd":   macd[len(macd)-2],
			"prev_signal": signal[len(signal)-2],
		},
	}
}
