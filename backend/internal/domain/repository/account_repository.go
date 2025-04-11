// internal/domain/repository/account_repository.go

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

// AccountRepository define a interface para qualquer banco de dados
type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	GetAll(ctx context.Context) ([]*entity.Account, error)
	UpdateByID(ctx context.Context, id uuid.UUID, jsonData []byte) (*entity.Account, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
