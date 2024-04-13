package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mailer/config"
	"time"
)

func New(ctx context.Context, params config.Mongo) *mongo.Database {
	connCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connCtx, options.Client().ApplyURI(params.Url))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(connCtx, nil); err != nil {
		panic(err)
	}

	context.AfterFunc(ctx, closeConn(client))
	return client.Database(params.DbName)
}

func closeConn(client *mongo.Client) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Printf("failed to close mongo connection: %v", err)
		}
	}
}
