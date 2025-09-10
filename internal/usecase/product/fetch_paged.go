package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type FetchPagedProductsRequest struct {
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	Sort       string `json:"sort"`
}

type FetchPagedProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type ProductResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FetchPagedProductsUseCase struct {
	productRepository domain.ProductRepositoryInterface
}

func NewFetchPagedProductsUseCase(productRepository domain.ProductRepositoryInterface) *FetchPagedProductsUseCase {
	return &FetchPagedProductsUseCase{
		productRepository: productRepository,
	}
}

func (uc *FetchPagedProductsUseCase) Execute(ctx context.Context, r FetchPagedProductsRequest) (FetchPagedProductsResponse, error) {
	products, err := uc.productRepository.FetchPaged(
		ctx,
		r.PageNumber,
		r.PageSize,
		r.Sort,
	)
	if err != nil {
		return FetchPagedProductsResponse{}, err
	}

	return FetchPagedProductsResponse{Products: mapProducts(products)}, nil
}

func mapProducts(products []domain.Product) []ProductResponse {
	outputs := make([]ProductResponse, len(products))
	for i, p := range products {
		outputs[i] = ProductResponse{
			Id:        p.Id,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return outputs
}
