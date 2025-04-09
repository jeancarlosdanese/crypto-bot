// internal/infra/mongo/asset_mongo.go

package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AssetMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Symbol    string             `bson:"symbol"`
	Quantity  float64            `bson:"quantity"`
	Category  string             `bson:"category"` // Ex: "reserve", "speculative", "high-risk"
	CreatedAt int64              `bson:"created_at"`
}
