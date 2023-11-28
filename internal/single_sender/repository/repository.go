package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate ifacemaker -f *.go -o ../repository.go -i Repository -s Repository -p single_sender -y "Controller describes methods, implemented by the repository package."
//go:generate mockgen -package mock -source ../repository.go -destination mock/repository_mock.go
type Repository struct {
	db *mongo.Collection
}

func NewSingleSenderRepo(db *mongo.Collection) *Repository {
	return &Repository{
		db: db,
	}
}
