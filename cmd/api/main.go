// cmd/api/main.go

package main

import (
	"log"
	"os"

	binanceApi "github.com/adshao/go-binance/v2"
	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/config"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/database"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/repository/postgres"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/services/binance"
)

func main() {
	logger.InitLogger()
	log.Println("🚀 Iniciando Robe de Crypto...")

	config.LoadEnv(".env")

	// PostgreSQL
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("Erro ao conectar no PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Repositórios
	accountRepo := postgres.NewAccountRepository(pool)
	botRepo := postgres.NewBotRepository(pool)
	positionRepo := postgres.NewPositionRepository(pool)
	executionRepo := postgres.NewExecutionLogRepository(pool)
	decisionRepo := postgres.NewDecisionLogRepository(pool)

	// Exchange Service (Binance)
	binanceClient := binanceApi.NewClient(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_API_SECRET"))
	exchangeService := binance.NewBinanceService(binanceClient)

	// Recupera todos os bots ativos
	// (Exemplo estático de um accountID — ideal seria iterar contas ou rodar por user autenticado)
	account, err := accountRepo.GetByEmail("jean@danese.com.br")
	if err != nil {
		log.Fatalf("Erro ao carregar conta: %v", err)
	}

	bots, err := botRepo.GetByAccountID(account.ID)
	if err != nil {
		log.Fatalf("Erro ao carregar bots: %v", err)
	}

	for _, bot := range bots {
		if !bot.Active {
			continue
		}

		go func(botInfo entity.Bot) {
			// Inicializa estratégia com repositórios injetados
			strategy := usecases.NewStrategyUseCase(*account, bot, exchangeService, decisionRepo, executionRepo, positionRepo, 240)

			// Restaura posição aberta (se houver)
			if pos, _ := positionRepo.Get(botInfo.ID); pos != nil {
				strategy.PositionQuantity = 1
				strategy.LastEntryPrice = pos.EntryPrice
				strategy.LastEntryTimestamp = pos.Timestamp
				log.Printf("🔁 [%s] Posição reaberta a %.2f", botInfo.Symbol, pos.EntryPrice)
			}

			// Inicia o monitoramento
			stream := binance.NewBinanceStreamService(strategy, exchangeService)
			stream.Start(botInfo.Symbol, botInfo.Interval)
		}(bot)
	}

	select {} // mantém o programa vivo
}
