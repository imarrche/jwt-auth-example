package app

import (
	"github.com/imarrche/jwt-auth-example/internal/service"
	"github.com/imarrche/jwt-auth-example/internal/store"
)

// Service is the app service implementation.
type Service struct {
	store store.Store
	auth  *authService
}

// NewService creates and returns a new service instance.
func NewService(s store.Store) *Service {
	return &Service{store: s}
}

// Auth returns authorization service.
func (s *Service) Auth() service.Auth {
	if s.auth == nil {
		s.auth = newAuthService(s.store)
	}

	return s.auth
}
