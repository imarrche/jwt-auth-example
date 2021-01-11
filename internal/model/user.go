package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User model represents a user.
type User struct {
	ID           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	FirstName    string `json:"first_name" db:"first_name"`
	SecondName   string `json:"second_name" db:"second_name"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-" db:"password_hash"`
}

// Validate validates user's fields.
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.FirstName, validation.Required, validation.Length(0, 50)),
		validation.Field(&u.SecondName, validation.Required, validation.Length(0, 50)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 256)),
	)
}
