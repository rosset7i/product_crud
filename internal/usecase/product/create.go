package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type CreateRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateResponse struct {
	Id uuid.UUID `json:"id"`
}

type CreateUseCase struct {
	productRepository domain.ProductRepositoryInterface
}

func NewCreateUseCase(productRepository domain.ProductRepositoryInterface) *CreateUseCase {
	return &CreateUseCase{
		productRepository: productRepository,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, r CreateRequest) (CreateResponse, error) {
	p, err := domain.NewProduct(r.Name, r.Price)
	if err != nil {
		return CreateResponse{}, nil
	}

	err = uc.productRepository.Create(ctx, p)
	if err != nil {
		return CreateResponse{}, nil
	}

	return CreateResponse{
		Id: p.Id,
	}, nil
}
