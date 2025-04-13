// internal/domain/repository/bot_repository.go

package repository

import (
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type BotRepository interface {
	Create(bot *entity.Bot) (*entity.BotWithStrategy, error)
	GetByID(id uuid.UUID) (*entity.BotWithStrategy, error)
	GetByAccountID(accountID uuid.UUID) ([]entity.BotWithStrategy, error)
	Update(bot *entity.Bot) (*entity.BotWithStrategy, error)
}
