// internal/domain/entity/asset.go

package entity

import (
	"time"
)

type Asset struct {
	ID        string    `bson:"id,omitempty"`
	Symbol    string    `bson:"symbol"`
	Quantity  float64   `bson:"quantity"`
	Category  string    `bson:"category"` // Ex: "reserve", "speculative", "high-risk"
	CreatedAt time.Time `bson:"created_at"`
}
