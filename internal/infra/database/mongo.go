// internal/infra/database/mongo.go

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() (*mongo.Client, error) {
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	db := os.Getenv("MONGO_DATABASE")

	uri := fmt.Sprintf("mongodb://%s:%s@%s", user, pass, host)
	clientOpts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Verificar conexão
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Printf("MongoDB conectado com sucesso ao host: %s (db: %s)", host, db)
	return client, nil
}
