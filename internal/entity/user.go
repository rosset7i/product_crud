package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	baseModel
	Name         string
	Email        string
	PasswordHash string
}

var (
	ErrEmailIsRequired = errors.New("email is required")
)

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		baseModel:    initEntity(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameIsRequired
	}
	if u.Email == "" {
		return ErrEmailIsRequired
	}

	return nil
}
