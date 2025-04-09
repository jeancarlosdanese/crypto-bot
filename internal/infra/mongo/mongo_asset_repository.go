// internal/infra/mongo/mongo_asset_repository.go

package mongo

import (
	"context"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAssetRepository struct {
	collection *mongo.Collection
}

func NewMongoAssetRepository(db *mongo.Database) repository.AssetRepository {
	return &mongoAssetRepository{
		collection: db.Collection("assets"),
	}
}

func (r *mongoAssetRepository) Save(ctx context.Context, asset *entity.Asset) error {
	assetMongo := toAssetMongo(asset)
	_, err := r.collection.InsertOne(ctx, assetMongo)
	return err
}

func (r *mongoAssetRepository) FindAll(ctx context.Context) ([]*entity.Asset, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assets []*entity.Asset
	for cursor.Next(ctx) {
		var assetMongo AssetMongo
		if err := cursor.Decode(&assetMongo); err != nil {
			return nil, err
		}
		assets = append(assets, toDomain(&assetMongo))
	}
	return assets, nil
}

func (r *mongoAssetRepository) FindBySymbol(ctx context.Context, symbol string) (*entity.Asset, error) {
	var assetMongo AssetMongo
	err := r.collection.FindOne(ctx, bson.M{"symbol": symbol}).Decode(&assetMongo)
	if err != nil {
		return nil, err
	}
	return toDomain(&assetMongo), nil
}
