package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	baseModel
	Name         string
	Email        string
	PasswordHash string
}

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		baseModel:    initEntity(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}, nil
}
