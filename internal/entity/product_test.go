package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Product", 10)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.Id)
	assert.Equal(t, "Product", product.Name)
	assert.Equal(t, float64(10), product.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 10)
	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, errNameIsRequired, err)
}

func TestProductWhenPriceIsZero(t *testing.T) {
	product, err := NewProduct("Product", 0)
	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, errPriceMustBeGreaterThanZero, err)
}

func TestProductValidate(t *testing.T) {
	product, err := NewProduct("Product", 10)
	assert.NotNil(t, product)
	assert.Nil(t, err)
	assert.Nil(t, product.Validate())
}
