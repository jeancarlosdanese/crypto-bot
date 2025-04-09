// internal/domain/repository/position_repository.go

package repository

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type PositionRepository interface {
	Save(position entity.OpenPosition) error
	Delete(symbol string) error
	Get(symbol string) (*entity.OpenPosition, error)
	GetAll() ([]entity.OpenPosition, error)
}
