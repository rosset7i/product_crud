package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/product_crud/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FetchByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`,
		email,
	).Scan(&u.Id, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(
		context.Background(),
		"INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}
