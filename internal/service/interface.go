package service

import "github.com/imarrche/jwt-auth-example/internal/model"

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

// Service is the interface all services must implement.
type Service interface {
	Auth() Auth
}

// Auth is the interface all authorization services must implement.
type Auth interface {
	SignUp(model.User) (model.User, error)
	SignIn(string, string) (string, string, error)
	ValidateJWT(string, string) (int, error)
	RefreshAccessJWT(string) (string, error)
}
