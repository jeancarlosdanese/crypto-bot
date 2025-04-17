// cmd/api/bot_runner.go

package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/factory"
	"github.com/jeancarlosdanese/crypto-bot/internal/runtime"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
	"github.com/jeancarlosdanese/crypto-bot/internal/services/binance"
)

func startBots(
	ctx context.Context,
	accountRepo repository.AccountRepository,
	botRepo repository.BotRepository,
	positionRepo repository.PositionRepository,
	decisionRepo repository.DecisionLogRepository,
	executionRepo repository.ExecutionLogRepository,
) {
	// 🔍 Por enquanto, buscamos só uma conta (exemplo fixo ou admin)
	accountID := uuid.MustParse(os.Getenv("ACCOUNT_ADMIN_ID"))
	account, err := accountRepo.GetByID(ctx, accountID)
	if err != nil {
		log.Fatalf("Erro ao carregar conta: %v", err)
	}

	// Exchange Service (Binance)
	exchangeFactory := factory.NewExchangeFactory()
	binanceExchangeService := exchangeFactory.NewExchangeService("binance", account)

	// 🔁 Start bots em paralelo
	binanceStreamFactory := func(strategy *usecases.StrategyUseCase) services.StreamService {
		return binance.NewBinanceStreamService(strategy, binanceExchangeService.(*binance.BinanceService))
	}

	bots, err := botRepo.GetByAccountID(account.ID)
	if err != nil {
		log.Fatalf("Erro ao carregar bots da conta: %v", err)
	}

	for _, bot := range bots {
		if !bot.Active {
			continue
		}

		go func(botInfo entity.BotWithStrategy) {
			strategyImpl, err := factory.NewStrategyByName(botInfo.StrategyName)
			if err != nil {
				log.Printf("❌ Estratégia desconhecida para bot %s: %v", botInfo.ID, err)
				return
			}

			strategy := usecases.NewStrategyUseCase(*account, botInfo, binanceExchangeService, strategyImpl, decisionRepo, executionRepo, positionRepo, 240)

			// Salvar no mapa global
			runtime.BotsMap.Lock()
			runtime.BotsMap.Items[botInfo.ID] = strategy
			runtime.BotsMap.Unlock()

			if pos, _ := positionRepo.Get(botInfo.ID); pos != nil {
				strategy.PositionQuantity = 1
				strategy.LastEntryPrice = pos.EntryPrice
				strategy.LastEntryTimestamp = pos.Timestamp
				log.Printf("🔁 [%s] Posição reaberta a %.2f", botInfo.Symbol, pos.EntryPrice)
			}

			stream := binanceStreamFactory(strategy)
			_ = stream.Start(botInfo.Symbol, botInfo.Interval)
		}(bot)
	}
}
