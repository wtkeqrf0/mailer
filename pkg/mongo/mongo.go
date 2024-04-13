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

	client, err := mongo.Connect(connCtx, options.Client().ApplyURI(params.URL))
	if err != nil {
		log.Fatalf("failed to init MongoDB conn: %v", err)
	}

	if err = client.Ping(connCtx, nil); err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	context.AfterFunc(ctx, closeConn(client))
	return client.Database(params.DbName)
}

func closeConn(client *mongo.Client) func() {
	return func() {
		c, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		if err := client.Disconnect(c); err != nil {
			log.Printf("failed to close mongo connection: %v", err)
		}
	}
}
