package entity

import "errors"

type Product struct {
	baseModel
	Name  string
	Price float64
}

var (
	ErrNameIsRequired             = errors.New("name is required")
	ErrPriceMustBeGreaterThanZero = errors.New("price must be greater than 0")
)

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		baseModel: initEntity(),
		Name:      name,
		Price:     price,
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price <= 0 {
		return ErrPriceMustBeGreaterThanZero
	}

	return nil
}
