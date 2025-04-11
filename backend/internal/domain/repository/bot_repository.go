// internal/domain/repository/bot_repository.go

package repository

import (
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type BotRepository interface {
	Create(bot *entity.Bot) (*entity.Bot, error)
	GetByID(id uuid.UUID) (*entity.Bot, error)
	GetByAccountID(accountID uuid.UUID) ([]entity.Bot, error)
	Update(bot *entity.Bot) (*entity.Bot, error)
}
