package models

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	Email             string
	Password          string
	EncryptedPassword string
}

var ErrUserExists = errors.New("user exists")

// Validate валидация логина и пароля
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Length(6, 100)),
	)
}
