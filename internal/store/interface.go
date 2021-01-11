package store

import "github.com/imarrche/jwt-auth-example/internal/model"

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

// Store is the interface all stores must implement.
type Store interface {
	Open() error
	Users() UserRepo
	Close() error
}

// UserRepo is the interface all user repositories must implement.
type UserRepo interface {
	GetAll() ([]model.User, error)
	Create(model.User) (model.User, error)
	GetByID(int) (model.User, error)
	GetByEmail(string) (model.User, error)
	Update(model.User) (model.User, error)
	DeleteByID(int) error
}
