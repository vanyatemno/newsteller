package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"newsteller/api/handlers"
	"newsteller/internal/cache"
	"newsteller/internal/config"
)

type Posts struct {
	handler *handlers.Post
}

func NewPosts(cfg *config.Config, c *mongo.Collection, cache *cache.PagesCache) *Posts {
	return &Posts{
		handler: handlers.NewPost(cfg, c, cache),
	}
}

func (p *Posts) SetRoutes(app *fiber.App) {
	postGroup := app.Group("/posts")
	postGroup.Post("/", p.handler.Create)
	postGroup.Delete("/:id", p.handler.Delete)
	postGroup.Put("/:id", p.handler.Update)
}
