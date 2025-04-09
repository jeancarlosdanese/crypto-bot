// internal/infra/mongo/mongo_decision_log_repository.go

package mongo

import (
	"context"
	"time"

	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDecisionLogRepository struct {
	collection *mongo.Collection
}

func NewDecisionLogRepository(db *mongo.Database) repository.DecisionLogRepository {
	col := db.Collection("decision_logs")

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "symbol", Value: 1}, {Key: "timestamp", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	col.Indexes().CreateOne(context.TODO(), index)

	return &mongoDecisionLogRepository{collection: col}
}

func (r *mongoDecisionLogRepository) Save(log entity.DecisionLog) error {
	log.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.TODO(), log)
	if err != nil {
		logger.Error("Erro ao salvar log de decisão", err, "symbol", log.Symbol, "timestamp", log.Timestamp)
	}
	return err
}
