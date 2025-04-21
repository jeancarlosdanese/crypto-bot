// internal/services/indicator_service.go

package services

import (
	"github.com/jeancarlosdanese/crypto-bot/internal/app/indicators"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type IndicatorService struct{}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{}
}

func (s *IndicatorService) GenerateSnapshot(
	candles []entity.Candle,
	emaPeriods []int,
	macdShort, macdLong, macdSignal int,
	rsiPeriod int,
	atrPeriod int,
	volatilityWindow int,
	bollingerPeriod int,
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

	emaMap := make(map[int]float64)
	for _, period := range emaPeriods {
		emaMap[period] = indicators.LastEMA(prices, period)
	}

	macd, signal, _ := indicators.MACD(prices, macdShort, macdLong, macdSignal)
	var macdVal, macdSignalVal float64
	if len(macd) > 0 {
		macdVal = macd[len(macd)-1]
	}
	if len(signal) > 0 {
		macdSignalVal = signal[len(signal)-1]
	}

	rsi := indicators.RSI(prices, rsiPeriod)
	atr := indicators.ATRFromCandles(candles[len(candles)-atrPeriod:])
	volatility := indicators.Volatility(prices[len(prices)-volatilityWindow:])
	bbUpper, bbLower := indicators.BollingerBands(prices, bollingerPeriod)
	bbWidth := bbUpper - bbLower
	avgVolume := indicators.AverageVolume(candles[len(candles)-11:])

	return &entity.IndicatorSnapshot{
		Timestamp:  candles[len(candles)-1].Time,
		Price:      prices[len(prices)-1],
		Volume:     volumes[len(volumes)-1],
		EMAs:       emaMap,
		MACD:       macdVal,
		MACDSignal: macdSignalVal,
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
