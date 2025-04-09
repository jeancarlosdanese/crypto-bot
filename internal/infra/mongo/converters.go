// internal/infra/mongo/converters.go

package mongo

import (
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toAssetMongo(asset *entity.Asset) *AssetMongo {
	objectID := primitive.NewObjectID()
	if asset.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(asset.ID); err == nil {
			objectID = oid
		}
	}

	return &AssetMongo{
		ID:        objectID,
		Symbol:    asset.Symbol,
		Quantity:  asset.Quantity,
		Category:  asset.Category,
		CreatedAt: asset.CreatedAt.Unix(),
	}
}

func toDomain(assetMongo *AssetMongo) *entity.Asset {
	return &entity.Asset{
		ID:        assetMongo.ID.Hex(),
		Symbol:    assetMongo.Symbol,
		Quantity:  assetMongo.Quantity,
		Category:  assetMongo.Category,
		CreatedAt: time.Unix(assetMongo.CreatedAt, 0),
	}
}
