package grpc

import "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"

type Repository struct {
	Recipe *Recipe
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	recipeService, err := NewRecipe(*cfg.RecipeService.Addr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Recipe: recipeService,
	}, nil
}

func (r *Repository) Stop() error {
	_ = r.Recipe.Conn.Close()
	return nil
}
