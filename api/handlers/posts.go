package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"newsteller/api/dto"
	"newsteller/internal/cache"
	"newsteller/internal/config"
	"newsteller/internal/models"
	"newsteller/internal/state"
	"time"
)

type Post struct {
	cfg   *config.Config
	state state.State[models.Post]
	cache *cache.PagesCache
}

func NewPost(cfg *config.Config, c *mongo.Collection, cache *cache.PagesCache) *Post {
	return &Post{
		cfg:   cfg,
		state: state.NewPostState(c),
		cache: cache,
	}
}

func (p *Post) Create(c *fiber.Ctx) error {
	var createPostDTO dto.PostDTO
	err := c.BodyParser(&createPostDTO)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(createPostDTO)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	err = p.state.Insert(c.Context(), &models.Post{
		Title:     createPostDTO.Title,
		Content:   createPostDTO.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	p.cache.Invalidate(cache.PostsUpdated)

	return c.SendStatus(fiber.StatusCreated)
}

func (p *Post) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	err := p.state.Delete(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	p.cache.Invalidate(cache.PostsUpdated)

	return c.SendStatus(fiber.StatusNoContent)
}

func (p *Post) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	var createPostDTO dto.PostDTO
	err := c.BodyParser(&createPostDTO)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(createPostDTO)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	err = p.state.Update(c.Context(), &models.Post{
		ID:        objectID,
		Title:     createPostDTO.Title,
		Content:   createPostDTO.Content,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	p.cache.Invalidate(cache.PostsUpdated)

	return c.SendStatus(fiber.StatusOK)
}
