package services

import "github.com/mephistolie/chefbook-backend-shopping-list/internal/config"

type Remote struct {
	Auth *Auth
}

func NewRemote(cfg *config.Config) (*Remote, error) {
	authService, err := NewAuth(*cfg.AuthService.Addr)
	if err != nil {
		return nil, err
	}

	return &Remote{
		Auth: authService,
	}, nil
}

func (s *Remote) Stop() error {
	return s.Auth.Conn.Close()
}
