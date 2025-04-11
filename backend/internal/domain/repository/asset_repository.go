// internal/domain/repository/asset_repository.go

package repository

import (
	"context"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type AssetRepository interface {
	Save(ctx context.Context, asset *entity.Asset) error
	FindAll(ctx context.Context) ([]*entity.Asset, error)
	FindBySymbol(ctx context.Context, symbol string) (*entity.Asset, error)
}
