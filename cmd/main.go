package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"newsteller/api/routes"
	"newsteller/internal/cache"
	"newsteller/internal/config"
	"newsteller/internal/db"
	"newsteller/internal/models"
	"os/signal"
	"syscall"
)

var (
	webApp *fiber.App
)

func main() {
	var cancel context.CancelFunc
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	cfg, err := config.Read()
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	client, err := db.Connect(ctx, cfg)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			zap.L().Fatal("failed to disconnect from database", zap.Error(err))
		}
	}()

	go createWebServer(cfg, client)

	for range ctx.Done() {
		_ = webApp.ShutdownWithContext(ctx)
		return
	}
}

func createWebServer(cfg *config.Config, client *mongo.Client) {
	webApp = fiber.New()
	setupWebServer(cfg, webApp, client)

	if err := webApp.Listen(":" + cfg.Port); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}

func setupWebServer(cfg *config.Config, app *fiber.App, client *mongo.Client) {
	pagesCache := cache.NewPagesCache()
	postsCollection := client.
		Database(cfg.Database.Name).
		Collection(models.Post{}.CollectionName())

	routes.New().InitializeRoutes(
		app,
		routes.NewPosts(cfg, postsCollection, pagesCache),
		routes.NewPages(cfg, postsCollection, pagesCache),
	)
}
