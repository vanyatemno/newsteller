package state

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"newsteller/internal/models"
	"newsteller/internal/repositories"
)

type Post struct {
	repo *repositories.Post
	// postsMap - map of posts with ID as a key
	postsMap map[string]*models.Post
}

func (p *Post) FindAll(ctx context.Context) ([]models.Post, error) {
	return p.repo.All(ctx)
}

func (p *Post) FindPaginated(ctx context.Context, q *repositories.PaginatedSearchQuery) ([]models.Post, int64, error) {
	return p.repo.FindPaginated(ctx, q)
}

func (p *Post) FindByID(ctx context.Context, id string) (*models.Post, error) {
	post, ok := p.postsMap[id]
	if ok {
		return post, nil
	}

	res, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.postsMap[id] = res

	return res, nil
}

func (p *Post) Insert(ctx context.Context, model *models.Post) error {
	res, err := p.repo.Create(ctx, model)
	if err != nil {
		return err
	}
	model.ID = *res

	return nil
}

func (p *Post) Delete(ctx context.Context, id string) error {
	delete(p.postsMap, id)
	return p.repo.Delete(ctx, id)
}

func (p *Post) Update(ctx context.Context, model *models.Post) error {
	p.postsMap[model.ID.Hex()] = model
	return p.repo.Update(ctx, model)
}

func NewPostState(c *mongo.Collection) State[models.Post] {
	return &Post{
		repo:     repositories.NewPostRepository(c),
		postsMap: make(map[string]*models.Post),
	}
}
