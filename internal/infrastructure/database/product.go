package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/product_crud/internal/domain"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) FetchPaged(ctx context.Context, pageNumber, pageSize int, sort string) ([]domain.Product, error) {
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	offset := (pageNumber - 1) * pageSize
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, name, price, created_at, updated_at
		FROM products
		ORDER BY name `+sort+`
		LIMIT $1 OFFSET $2`,
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *ProductRepository) FetchById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, name, price, created_at, updated_at
		FROM products
		WHERE id = $1`,
		id,
	).Scan(&p.Id, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.db.Exec(
		context.Background(),
		"INSERT INTO products (id, name, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		product.Id,
		product.Name,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	cmd, err := r.db.Exec(
		context.Background(),
		"UPDATE products SET (name, price, updated_at) = ($1, $2, $3) WHERE id = $4",
		product.Name,
		product.Price,
		product.UpdatedAt,
		product.Id,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("product not found")
	}

	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.db.Exec(
		context.Background(),
		"DELETE FROM products WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("product not found")
	}

	return nil
}
