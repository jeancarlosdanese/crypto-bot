// internal/app/indicators/indicators.go

package indicators

import (
	"math"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

// Volatility calcula o desvio padrão dos preços como uma medida de volatilidade.
func Volatility(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	mean := SMA(prices)
	sumSquaredDiffs := 0.0
	for _, price := range prices {
		diff := price - mean
		sumSquaredDiffs += diff * diff
	}
	variance := sumSquaredDiffs / float64(len(prices))
	// Desvio padrão relativo à média, em porcentagem
	return (math.Sqrt(variance) / mean) * 100
}

// SMA calcula a Média Móvel Simples para um slice de preços.
func SMA(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	return sum / float64(len(prices))
}

// MovingAverage calcula a média móvel simples para os últimos 'period' valores da janela (usando preços de fechamento).
func MovingAverage(prices []float64, period int) float64 {
	if len(prices) < period {
		return SMA(prices)
	}
	slice := prices[len(prices)-period:]
	return SMA(slice)
}

// EMA calcula a Média Móvel Exponencial para um slice de preços.
// O cálculo utiliza o primeiro valor como ponto de partida.
func EMA(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	alpha := 2.0 / (float64(len(prices)) + 1.0)
	ema := prices[0]
	for i := 1; i < len(prices); i++ {
		ema = alpha*prices[i] + (1-alpha)*ema
	}
	return ema
}

// ATRPercent calcula o ATR como percentual do preço de fechamento do candle mais recente.
func ATRPercent(candles []entity.Candle) float64 {
	atrAbsolute := ATRFromCandles(candles)
	if len(candles) == 0 {
		return 0
	}
	lastClose := candles[len(candles)-1].Close
	return (atrAbsolute / lastClose) * 100
}

// ATRFromCandles calcula o Average True Range (ATR) a partir de um slice de Candle.
// Essa função extrai os valores de high, low e close de cada candle e utiliza a função ATR já existente.
func ATRFromCandles(candles []entity.Candle) float64 {
	if len(candles) < 2 {
		return 0
	}
	var highs, lows, closes []float64
	for _, c := range candles {
		highs = append(highs, c.High)
		lows = append(lows, c.Low)
		closes = append(closes, c.Close)
	}
	return ATR(highs, lows, closes)
}

// ATR (Average True Range) é uma medida de volatilidade que considera
// o maior valor entre a variação do período atual e a diferença com o fechamento anterior.
// É necessário fornecer slices de máximos, mínimos e fechamentos.
func ATR(highs, lows, closes []float64) float64 {
	n := len(closes)
	if n < 2 {
		return 0
	}
	trs := make([]float64, n-1)
	for i := 1; i < n; i++ {
		currentHigh := highs[i]
		currentLow := lows[i]
		prevClose := closes[i-1]

		tr1 := currentHigh - currentLow
		tr2 := math.Abs(currentHigh - prevClose)
		tr3 := math.Abs(currentLow - prevClose)

		trs[i-1] = math.Max(tr1, math.Max(tr2, tr3))
	}
	return SMA(trs)
}

// RSI calcula o Relative Strength Index para um slice de preços e um período.
// Se não houver dados suficientes, retorna 0.
func RSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 0
	}

	// Calcula os ganhos e perdas iniciais
	gain := 0.0
	loss := 0.0
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gain += change
		} else {
			loss -= change // transforma em valor absoluto
		}
	}

	avgGain := gain / float64(period)
	avgLoss := loss / float64(period)

	// Se a média das perdas for zero, retorna 100 (situação de forte momentum de alta)
	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

// EMASeries calcula a série da Média Móvel Exponencial para os preços fornecidos em um dado período.
func EMASeries(prices []float64, period int) []float64 {
	if len(prices) == 0 {
		return nil
	}
	result := make([]float64, len(prices))
	alpha := 2.0 / (float64(period) + 1.0)
	result[0] = prices[0]
	for i := 1; i < len(prices); i++ {
		result[i] = alpha*prices[i] + (1-alpha)*result[i-1]
	}
	return result
}

// MACD calcula o MACD, a linha de sinal e o histograma para uma série de preços.
// shortPeriod e longPeriod são os períodos para as EMAs de curto e longo prazo, respectivamente,
// e signalPeriod é o período para calcular a linha de sinal a partir do MACD.
func MACD(prices []float64, shortPeriod, longPeriod, signalPeriod int) (macdLine, signalLine, histogram []float64) {
	if len(prices) < longPeriod {
		return nil, nil, nil
	}

	// Calcula a EMA para os períodos curto e longo.
	emaShort := EMASeries(prices, shortPeriod)
	emaLong := EMASeries(prices, longPeriod)

	macdLine = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		macdLine[i] = emaShort[i] - emaLong[i]
	}

	// Calcula a linha de sinal a partir do MACD.
	signalLine = EMASeries(macdLine, signalPeriod)

	// O histograma é a diferença entre o MACD e a linha de sinal.
	histogram = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	return macdLine, signalLine, histogram
}
