package domain

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
	errEmailIsRequired = errors.New("email is required")
)

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		baseModel:    initEntity(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err = u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) Validate() error {
	switch {
	case u.Name == "":
		return errNameIsRequired
	case u.Email == "":
		return errEmailIsRequired
	}

	return nil
}
