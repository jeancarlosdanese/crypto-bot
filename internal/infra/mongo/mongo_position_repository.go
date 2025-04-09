// internal/infra/mongo/mongo_position_repository.go

package mongo

import (
	"context"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoPositionRepository struct {
	collection *mongo.Collection
}

func NewPositionRepository(db *mongo.Database) repository.PositionRepository {
	col := db.Collection("positions")

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "symbol", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	col.Indexes().CreateOne(context.TODO(), index)

	return &mongoPositionRepository{collection: col}
}

func (r *mongoPositionRepository) Save(position entity.OpenPosition) error {
	filter := bson.M{"symbol": position.Symbol}
	update := bson.M{"$set": position}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (r *mongoPositionRepository) Delete(symbol string) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"symbol": symbol})
	return err
}

func (r *mongoPositionRepository) Get(symbol string) (*entity.OpenPosition, error) {
	var pos entity.OpenPosition
	err := r.collection.FindOne(context.TODO(), bson.M{"symbol": symbol}).Decode(&pos)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &pos, err
}

func (r *mongoPositionRepository) GetAll() ([]entity.OpenPosition, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var positions []entity.OpenPosition
	if err = cursor.All(context.TODO(), &positions); err != nil {
		return nil, err
	}
	return positions, nil
}
