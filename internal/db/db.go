package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"newsteller/internal/config"
	"newsteller/internal/models"
)

func Connect(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	co := options.Client().ApplyURI(cfg.DNS)
	co.Auth = &options.Credential{
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
	}
	client, err := mongo.Connect(ctx, co)
	if err != nil {
		zap.L().Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	err = migrate(
		ctx,
		cfg,
		client,
		models.Post{},
	)
	if err != nil {
		zap.L().Error("failed to migrate database", zap.Error(err))
	}

	zap.L().Info("successfully connected to database and run migrations", zap.Any("config", cfg))

	return client, nil
}

func migrate(ctx context.Context, cfg *config.Config, client *mongo.Client, models ...models.Model) error {
	for i := range models {
		err := models[i].Migrate(ctx, client.Database(cfg.Database.Name))
		if err != nil {
			return err
		}
	}

	return nil
}
