package grpc

import "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"

type Repository struct {
	Profile *Profile
	Recipe  *Recipe
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	profileService, err := NewProfile(*cfg.ProfileService.Addr)
	if err != nil {
		return nil, err
	}
	recipeService, err := NewRecipe(*cfg.RecipeService.Addr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Profile: profileService,
		Recipe:  recipeService,
	}, nil
}

func (r *Repository) Stop() error {
	_ = r.Profile.Conn.Close()
	_ = r.Recipe.Conn.Close()
	return nil
}
