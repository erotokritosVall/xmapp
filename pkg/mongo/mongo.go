package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	Uri string `envconfig:"MONGO_URI"`
	Db  string `envconfig:"MONGO_DB"`
}

func New(ctx context.Context, cfg *Configuration) (*mongo.Database, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		return nil, fmt.Errorf("mongo failed to connect: %+v", err)
	}

	if err := conn.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo failed to ping: %+v", err)
	}

	return conn.Database(cfg.Db), nil
}
