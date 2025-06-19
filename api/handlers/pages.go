package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"newsteller/internal/cache"
	"newsteller/internal/config"
	"newsteller/internal/models"
	"newsteller/internal/repositories"
	"newsteller/internal/state"
	"newsteller/internal/templates"
)

type Page struct {
	cfg   *config.Config
	state state.State[models.Post]
	cache *cache.PagesCache
}

func NewPage(cfg *config.Config, c *mongo.Collection, cache *cache.PagesCache) *Page {
	return &Page{
		cfg:   cfg,
		state: state.NewPostState(c),
		cache: cache,
	}
}

// GET /home
func (p *Page) GetHomePage(c *fiber.Ctx) error {
	res, _, err := p.state.FindPaginated(c.Context(), &repositories.PaginatedSearchQuery{
		Page:  1,
		Limit: p.cfg.PostsPerPage,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	html, err := templates.NewMain(res).GeneratePage()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET /posts/search
func (p *Page) FindPaginated(c *fiber.Ctx) error {
	query, err := p.validatePaginationQuery(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, total, err := p.state.FindPaginated(c.Context(), query)
	if err != nil {
		zap.L().Error("could not get all posts", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	html, err := templates.
		NewHome(
			res,
			query.Page,
			int(total),
			p.cfg.PostsPerPage,
		).
		GeneratePage()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET /posts
func (p *Page) FindPostsList(c *fiber.Ctx) error {
	query, err := p.validatePaginationQuery(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, total, err := p.state.FindPaginated(c.Context(), query)
	if err != nil {
		zap.L().Error("could not get all posts", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	html, err := templates.
		NewList(
			res,
			query.Page,
			int(total),
			p.cfg.PostsPerPage,
		).
		GeneratePage()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET posts/:id
func (p *Page) FindPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	zap.L().Info("Getting post", zap.String("id", id))
	post, err := p.state.FindByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	html, err := templates.RenderSinglePost(post)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET /posts/create
func (p *Page) GetCreatePage(c *fiber.Ctx) error {
	html, err := templates.NewCreatePage().GeneratePage()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET /posts/moderation
func (p *Page) GetModerationPage(c *fiber.Ctx) error {
	query, err := p.validatePaginationQuery(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	posts, total, err := p.state.FindPaginated(c.Context(), query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	html, err := templates.
		NewModeration(
			posts,
			query.Page,
			query.Limit,
			int(total),
		).
		GeneratePage()

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

// GET /posts/:id/edit
func (p *Page) GetEditPage(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := p.state.FindByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "post with provided id does not exist")
	}

	html, err := templates.
		NewEdit(post).
		GeneratePage()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

func (p *Page) validatePaginationQuery(c *fiber.Ctx) (*repositories.PaginatedSearchQuery, error) {
	var query repositories.PaginatedSearchQuery
	err := c.QueryParser(&query)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = p.cfg.PostsPerPage
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(query)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return &query, nil
}
