package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	CollectionName() string
	Migrate(ctx context.Context, db *mongo.Database) error
}

//type Migratable interface {
//	Migrate(ctx context.Context, db *mongo.Database) error
//}
