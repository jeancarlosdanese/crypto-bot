// internal/report/summary.go

package reporter

import (
	"fmt"
	"os"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
)

func PrintExecutionSummary(executionLogRepo repository.ExecutionLogRepository) {
	executions, err := executionLogRepo.GetAll()
	if err != nil {
		logger.Error("Erro ao buscar execution_logs", err)
		return
	}

	type stats struct {
		trades   int
		win      int
		loss     int
		profit   float64
		roiTotal float64
	}

	summary := make(map[string]*stats)

	for _, e := range executions {
		s := summary[e.Symbol]
		if s == nil {
			s = &stats{}
			summary[e.Symbol] = s
		}
		s.trades++
		s.profit += e.Profit
		s.roiTotal += e.ROIPct
		if e.Profit > 0 {
			s.win++
		} else {
			s.loss++
		}
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("tmp/report_%s.log", timestamp)
	file, err := os.Create(fileName)
	if err != nil {
		logger.Error("Erro ao criar arquivo de log de resumo", err)
		return
	}
	defer file.Close()

	for symbol, s := range summary {
		winRate := float64(s.win) / float64(s.trades) * 100
		avgROI := s.roiTotal / float64(s.trades)

		logLine := fmt.Sprintf("📊 [%s] Trades: %d | Wins: %d | Losses: %d | Total Profit: %.2f | Avg ROI: %.2f%% | Win Rate: %.2f%%",
			symbol, s.trades, s.win, s.loss, s.profit, avgROI, winRate,
		)

		logger.Info("📊 Resumo de performance",
			"symbol", symbol,
			"trades", s.trades,
			"wins", s.win,
			"losses", s.loss,
			"total_profit", s.profit,
			"avg_roi_pct", avgROI,
			"win_rate_pct", winRate,
		)

		file.WriteString(logLine)
	}

	logger.Info("✅ Resumo salvo no arquivo", "file", fileName)
}
