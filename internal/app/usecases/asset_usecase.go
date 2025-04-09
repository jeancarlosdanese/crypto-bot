// internal/app/usecases/asset_usecase.go

package usecases

import (
	"context"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
)

type AssetUseCase struct {
	repo repository.AssetRepository
}

func NewAssetUseCase(repo repository.AssetRepository) *AssetUseCase {
	return &AssetUseCase{repo: repo}
}

func (s *AssetUseCase) RegisterAsset(ctx context.Context, symbol string, quantity float64, category string) error {
	asset := &entity.Asset{
		Symbol:   symbol,
		Quantity: quantity,
		Category: category,
	}
	return s.repo.Save(ctx, asset)
}

func (s *AssetUseCase) ListAssets(ctx context.Context) ([]*entity.Asset, error) {
	return s.repo.FindAll(ctx)
}

func (uc *AssetUseCase) GetReserveAssets(ctx context.Context) ([]*entity.Asset, error) {
	allAssets, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var reserves []*entity.Asset
	for _, asset := range allAssets {
		if asset.Category == "reserve" {
			reserves = append(reserves, asset)
		}
	}

	return reserves, nil
}
