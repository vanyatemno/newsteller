package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"go.mongodb.org/mongo-driver/mongo"
	"newsteller/api/handlers"
	"newsteller/internal/cache"
	"newsteller/internal/config"
)

type Pages struct {
	handler *handlers.Page
	cache   *cache.PagesCache
}

func NewPages(cfg *config.Config, c *mongo.Collection, cache *cache.PagesCache) *Pages {
	return &Pages{
		handler: handlers.NewPage(cfg, c, cache),
		cache:   cache,
	}
}

func (p *Pages) SetRoutes(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		if html, ok := p.cache.Get(c.Request().URI().String()); ok {
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
			return c.SendString(html)
		}
		err := c.Next()
		if err != nil {
			return err
		}
		if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			p.cache.Set(c.Request().URI().String(), string(c.Response().Body()))
		}

		return nil
	})
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{"/": "home"},
	}))

	app.Get("/home", p.handler.GetHomePage)

	postsGroup := app.Group("/posts")
	postsGroup.Get("/", p.handler.FindPostsList)
	postsGroup.Get("/search", p.handler.FindPaginated)
	postsGroup.Get("/create", p.handler.GetCreatePage)
	postsGroup.Get("/edit", p.handler.GetModerationPage)
	postsGroup.Get("/:id", p.handler.FindPostByID)
	postsGroup.Get("/:id/edit", p.handler.GetEditPage)
}
