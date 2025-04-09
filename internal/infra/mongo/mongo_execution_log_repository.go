// internal/infra/mongo/mongo_execution_log_repository.go

package mongo

import (
	"context"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoExecutionLogRepository struct {
	collection *mongo.Collection
}

func NewExecutionLogRepository(db *mongo.Database) repository.ExecutionLogRepository {
	col := db.Collection("execution_logs")

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "symbol", Value: 1}, {Key: "entry.timestamp", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	col.Indexes().CreateOne(context.TODO(), index)

	return &mongoExecutionLogRepository{collection: col}
}

func (r *mongoExecutionLogRepository) Save(log entity.ExecutionLog) error {
	log.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.TODO(), log)
	return err
}

func (r *mongoExecutionLogRepository) GetAll() ([]entity.ExecutionLog, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var logs []entity.ExecutionLog
	if err := cursor.All(context.TODO(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
