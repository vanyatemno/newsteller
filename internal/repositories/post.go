package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"newsteller/internal/models"
)

type Post struct {
	c *mongo.Collection
}

func NewPostRepository(collection *mongo.Collection) *Post {
	return &Post{c: collection}
}

func (p *Post) All(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	movieCursor, err := p.c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer movieCursor.Close(ctx)

	err = movieCursor.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (p *Post) FindByID(ctx context.Context, id string) (*models.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		zap.L().Error("could not convert string ID to primitive.ObjectID", zap.Error(err))
		return nil, err
	}

	var post models.Post
	err = p.c.FindOne(ctx, bson.D{{"_id", objectID}}).Decode(&post)
	if err != nil {
		zap.L().Error("could not find post by id", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return &post, nil
}

func (p *Post) Create(ctx context.Context, post *models.Post) (*primitive.ObjectID, error) {
	res, err := p.c.InsertOne(ctx, post)
	if err != nil {
		zap.L().Error("could not insert post", zap.Error(err))
		return nil, err
	}

	// Type assert the InsertedID to primitive.ObjectID
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		zap.L().Error("failed to convert InsertedID to ObjectID")
		return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
	}

	return &insertedID, nil
}

func (p *Post) FindPaginated(
	ctx context.Context,
	query *PaginatedSearchQuery,
) ([]models.Post, int64, error) {
	filter := bson.M{}
	if query.Keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query.Keyword, "$options": "i"}},
			{"content": bson.M{"$regex": query.Keyword, "$options": "i"}},
		}
	}

	skip := (query.Page - 1) * query.Limit
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.Limit)).
		SetSort(bson.D{{"created_at", -1}})

	cursor, err := p.c.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	total, err := p.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (p *Post) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		zap.L().Error("could not convert string ID to primitive.ObjectID", zap.Error(err))
		return err
	}

	_, err = p.c.DeleteOne(ctx, bson.D{{"_id", objectID}})
	if err != nil {
		zap.L().Error("could not delete post", zap.String("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (p *Post) Update(ctx context.Context, post *models.Post) error {
	objectID, err := primitive.ObjectIDFromHex(post.ID.Hex())
	if err != nil {
		zap.L().Error("could not convert string ID to primitive.ObjectID", zap.Error(err))
		return err
	}

	filter := bson.D{{"_id", objectID}}

	update := bson.D{
		{"$set", bson.D{
			{"title", post.Title},
			{"content", post.Content},
			{"updated_at", post.UpdatedAt},
		}},
	}

	result, err := p.c.UpdateOne(ctx, filter, update)
	if err != nil {
		zap.L().Error("could not update post", zap.String("id", post.ID.Hex()), zap.Error(err))
		return err
	}

	// check if any document was actually updated
	if result.MatchedCount == 0 {
		zap.L().Warn("no post found with given ID", zap.String("id", post.ID.Hex()))
		return fmt.Errorf("no post found with ID: %s", post.ID.Hex())
	}

	zap.L().Info("post updated successfully", zap.String("id", post.ID.Hex()))
	return nil
}
