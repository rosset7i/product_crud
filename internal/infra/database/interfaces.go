package database

import (
	"github.com/google/uuid"
	"github.com/rosset7i/zippy/internal/entity"
)

type UserInterface interface {
	Create(user *entity.User) error
	FetchByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FetchPaged(page, limit int, sortedBy string) ([]entity.Product, error)
	FetchById(id uuid.UUID) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uuid.UUID) error
}
