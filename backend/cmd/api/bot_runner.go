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
	"github.com/jeancarlosdanese/crypto-bot/internal/runtime"
	"github.com/jeancarlosdanese/crypto-bot/internal/services"
)

func startBots(
	ctx context.Context,
	accountRepo repository.AccountRepository,
	botRepo repository.BotRepository,
	positionRepo repository.PositionRepository,
	exchangeService services.ExchangeService,
	streamFactory func(*usecases.StrategyUseCase) services.StreamService,
	decisionRepo repository.DecisionLogRepository,
	executionRepo repository.ExecutionLogRepository,
) {
	// 🔍 Por enquanto, buscamos só uma conta (exemplo fixo ou admin)
	accountID := uuid.MustParse(os.Getenv("ACCOUNT_ADMIN_ID"))
	account, err := accountRepo.GetByID(ctx, accountID)
	if err != nil {
		log.Fatalf("Erro ao carregar conta: %v", err)
	}

	bots, err := botRepo.GetByAccountID(account.ID)
	if err != nil {
		log.Fatalf("Erro ao carregar bots da conta: %v", err)
	}

	for _, bot := range bots {
		if !bot.Active {
			continue
		}

		go func(botInfo entity.Bot) {
			strategy := usecases.NewStrategyUseCase(*account, botInfo, exchangeService, decisionRepo, executionRepo, positionRepo, 240)

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

			stream := streamFactory(strategy)
			_ = stream.Start(botInfo.Symbol, botInfo.Interval)
		}(bot)
	}
}
