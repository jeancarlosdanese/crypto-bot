// cmd/api/main.go

package main

import (
	"log"
	"os"

	"github.com/adshao/go-binance/v2"
	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/config"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/database"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/mongo"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	binanceService "github.com/jeancarlosdanese/crypto-bot/internal/services/binance"
)

func main() {
	logger.InitLogger()
	log.Println("Iniciando projeto...")

	config.LoadEnv(".env")

	mongoClient, err := database.NewMongoClient()
	if err != nil {
		log.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(nil)

	db := mongoClient.Database(os.Getenv("MONGO_DATABASE"))

	decisionLogRepo := mongo.NewDecisionLogRepository(db)
	executionLogRepo := mongo.NewExecutionLogRepository(db)
	positionRepo := mongo.NewPositionRepository(db)

	binanceClient := binance.NewClient(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_API_SECRET"))
	exchange := binanceService.NewBinanceService(binanceClient)

	windowSize := 240

	pairs := map[string]string{
		"btcusdt": "1m",
		"ethusdt": "1m",
		"solusdt": "1m",
	}

	for symbol, interval := range pairs {
		go func(sym, intv string) {
			strategy := usecases.NewStrategyUseCase(exchange, decisionLogRepo, executionLogRepo, positionRepo, windowSize)

			// Tenta restaurar posição salva (resiliência)
			if pos, _ := positionRepo.Get(sym); pos != nil {
				strategy.LastEntryPrice = pos.EntryPrice
				strategy.LastEntryTimestamp = pos.Timestamp
				strategy.PositionQuantity = 1
				log.Printf("🔁 Posição reaberta: %s @ %.2f", sym, pos.EntryPrice)
			}

			stream := binanceService.NewBinanceStreamService(strategy, exchange)
			stream.Start(sym, intv)
		}(symbol, interval)
	}

	select {} // mantém os goroutines vivos
}
