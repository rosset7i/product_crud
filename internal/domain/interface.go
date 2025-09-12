package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	FetchByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type ProductRepository interface {
	FetchPaged(ctx context.Context, pageNumber, pageSize int, sort string) ([]*Product, error)
	FetchById(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}
