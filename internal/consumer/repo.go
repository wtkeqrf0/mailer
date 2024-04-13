package consumer

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate ifacemaker -f *.go -o repo_if.go -i Repository -s repo -p consumer
type repo struct {
	db *mongo.Collection
}

func NewConsumerRepo(db *mongo.Collection) *repo {
	return &repo{
		db: db,
	}
}

// GetTemplateByName from the db, and save it into given part, if it can be found by name.
func (r *repo) GetTemplateByName(email *Email) error {
	bsonFilter, err := bson.Marshal(email.Settings)
	if err != nil {
		return err
	}

	if err = r.db.FindOne(context.Background(), bsonFilter).Decode(email); err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
