package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mailer/internal/consumer"
)

//go:generate ifacemaker -f *.go -o ../repository.go -i Repository -s Repository -p consumer -y "Controller describes methods, implemented by the repository package."
//go:generate mockgen -package mock -source ../repository.go -destination mock/repository_mock.go
type Repository struct {
	db *mongo.Collection
}

func NewConsumerRepo(db *mongo.Collection) *Repository {
	return &Repository{
		db: db,
	}
}

// GetTemplateByName from the db, and save it into given part, if it can be found by name.
func (r *Repository) GetTemplateByName(email *consumer.Email) error {
	bsonFilter, err := bson.Marshal(email.Settings)
	if err != nil {
		return err
	}

	if err = r.db.FindOne(context.Background(), bsonFilter).Decode(email); err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
