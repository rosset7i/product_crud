package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rosset7i/zippy/internal/entity"
)

type Product struct {
	DB *sql.DB
}

func (p *Product) Create(product *entity.Product) error {
	_, err := p.DB.Exec(
		"INSERT INTO products (id, name, price, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		product.Id,
		product.Name,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (p *Product) FetchPaged(pageNumber, pageSize int, _sortedBy string) ([]entity.Product, error) {
	var products []entity.Product
	offset := (pageNumber - 1) * pageSize
	rows, err := p.DB.Query("SELECT id, name, price, created_at, updated_at FROM products LIMIT $2 OFFSET $3", pageSize, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *Product) FetchById(id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	rows := p.DB.QueryRow("SELECT id, name, price, created_at, updated_at FROM products WHERE id = $1", id)
	err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt)

	return &product, err
}

func (p *Product) Update(product *entity.Product) error {
	product, err := p.FetchById(product.Id)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(
		"UPDATE products SET (name, price, updated_at) = ($1, $2, $3) WHERE id = $4",
		product.Name,
		product.Price,
		product.UpdatedAt,
		product.Id,
	)
	return err
}

func (p *Product) Delete(product *entity.Product) error {
	product, err := p.FetchById(product.Id)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec("DELETE FROM products WHERE id = $1", product.Id)
	return err
}
