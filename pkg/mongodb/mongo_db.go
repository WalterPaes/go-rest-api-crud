package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(parentCtx context.Context, timeout int, uri string) *mongo.Client {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().SetTimeout(time.Duration(timeout)*time.Second), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}
