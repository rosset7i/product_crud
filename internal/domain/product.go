package domain

import "errors"

type Product struct {
	baseModel
	Name  string
	Price float64
}

var (
	errNameIsRequired             = errors.New("name is required")
	errPriceMustBeGreaterThanZero = errors.New("price must be greater than 0")
)

func NewProduct(name string, price float64) (*Product, error) {
	p := &Product{
		baseModel: initEntity(),
		Name:      name,
		Price:     price,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Product) Validate() error {
	switch {
	case p.Name == "":
		return errNameIsRequired
	case p.Price <= 0:
		return errPriceMustBeGreaterThanZero
	}

	return nil
}
