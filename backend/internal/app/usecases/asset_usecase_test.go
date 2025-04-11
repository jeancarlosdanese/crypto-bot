// test/usecases/asset_usecase_test.go

package usecases_test

import (
	"context"
	"testing"

	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetReserveAssets(t *testing.T) {
	mockRepo := &mocks.MockAssetRepository{
		Assets: []*entity.Asset{
			{Symbol: "BTC", Category: "reserve"},
			{Symbol: "ETH", Category: "reserve"},
			{Symbol: "DOGE", Category: "speculative"},
		},
	}

	usecase := usecases.NewAssetUseCase(mockRepo)

	reserves, err := usecase.GetReserveAssets(context.Background())

	assert.NoError(t, err)
	assert.Len(t, reserves, 2)
	assert.Equal(t, "BTC", reserves[0].Symbol)
	assert.Equal(t, "ETH", reserves[1].Symbol)
}
