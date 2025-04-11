// test/mocks/mock_asset_repository.go

package mocks

import (
	"context"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
)

type MockAssetRepository struct {
	Assets []*entity.Asset
	Err    error
}

func (m *MockAssetRepository) FindAll(ctx context.Context) ([]*entity.Asset, error) {
	return m.Assets, m.Err
}

func (m *MockAssetRepository) FindBySymbol(ctx context.Context, symbol string) (*entity.Asset, error) {
	for _, a := range m.Assets {
		if a.Symbol == symbol {
			return a, nil
		}
	}
	return nil, nil
}

func (m *MockAssetRepository) Save(ctx context.Context, asset *entity.Asset) error {
	m.Assets = append(m.Assets, asset)
	return nil
}

var _ repository.AssetRepository = (*MockAssetRepository)(nil)
