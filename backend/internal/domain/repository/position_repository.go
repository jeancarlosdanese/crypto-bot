// internal/domain/repository/position_repository.go

package repository

import (
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type PositionRepository interface {
	Save(position entity.OpenPosition) error
	Delete(botID uuid.UUID) error
	Get(botID uuid.UUID) (*entity.OpenPosition, error)
	GetAll() ([]entity.OpenPosition, error)
}
