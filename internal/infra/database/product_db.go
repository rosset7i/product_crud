package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rosset7i/zippy/internal/entity"
)

type Product struct {
	DB *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{
		DB: db,
	}
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

func (p *Product) FetchPaged(pageNumber, pageSize int, sort string) ([]entity.Product, error) {
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	var products []entity.Product
	offset := (pageNumber - 1) * pageSize
	rows, err := p.DB.Query(fmt.Sprintf("SELECT id, name, price, created_at, updated_at FROM products ORDER BY created_at %v LIMIT $1 OFFSET $2", sort), pageSize, offset)
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
	_, err := p.FetchById(product.Id)
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

func (p *Product) Delete(id uuid.UUID) error {
	product, err := p.FetchById(id)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec("DELETE FROM products WHERE id = $1", product.Id)
	return err
}
