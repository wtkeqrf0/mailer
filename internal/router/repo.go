package router

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mailer/pkg/mail"
)

//go:generate ifacemaker -f *.go -o repo_if.go -i Repository -s repo -p router
type repo struct {
	db *mongo.Collection
}

func NewRepo(db *mongo.Collection) Repository {
	return &repo{
		db: db,
	}
}

// GetTemplateByName from the db, and save it into given part, if it can be found by name.
func (r *repo) GetTemplateByName(email *mail.Parsable) error {
	if email.Settings == nil {
		return nil
	}

	bsonFilter, err := bson.Marshal(email.Settings)
	if err != nil {
		return err
	}

	if err = r.db.FindOne(context.Background(), bsonFilter).Decode(email); err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
