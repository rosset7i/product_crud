package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type FetchByIdRequest struct {
	Id uuid.UUID `json:"id"`
}

type FetchByIdResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FetchByIdUseCase struct {
	productRepository domain.ProductRepositoryInterface
}

func NewFetchByIdUseCase(productRepository domain.ProductRepositoryInterface) *FetchByIdUseCase {
	return &FetchByIdUseCase{
		productRepository: productRepository,
	}
}

func (uc *FetchByIdUseCase) Execute(ctx context.Context, r FetchByIdRequest) (FetchByIdResponse, error) {
	p, err := uc.productRepository.FetchById(ctx, r.Id)
	if err != nil {
		return FetchByIdResponse{}, nil
	}

	return FetchByIdResponse{
		Id:        p.Id,
		Name:      p.Name,
		Price:     p.Price,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}
