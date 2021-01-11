package store

import "errors"

var (
	ErrUsernameIsTaken = errors.New("user with this username already exists")
	ErrEmailIsTaken    = errors.New("user with this email already exists")
)
