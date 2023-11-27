package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go-db-adapter/main/app/domain"
)

// Storage is the adapter interface that Clients will use.
type Storage interface {
	Close() error
	FindAll() ([]domain.Person, error)
	GetByID(id string) (domain.Person, error)
	Insert(data domain.Person) error
	GetSession() (mongo.Session, error)
	WithTransaction(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error
}
