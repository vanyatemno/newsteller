package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title,omitempty"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

func (Post) CollectionName() string {
	return "posts"
}

func (p Post) Migrate(ctx context.Context, db *mongo.Database) error {
	err := db.CreateCollection(ctx, p.CollectionName())
	if err != nil {
		return err
	}

	_, err = db.Collection(p.CollectionName()).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"title", "text"},
			{"content", "text"},
		},
	})

	return err
}
