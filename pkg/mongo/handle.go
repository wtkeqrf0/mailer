package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func ToGRPC(err error) error {
	switch err {
	case mongo.ErrNoDocuments, mongo.ErrEmptySlice:
		return status.Error(codes.NotFound, err.Error())
	}
	log.Println(err.Error())
	return err
}
