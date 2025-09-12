package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type UpdateRequest struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

type UpdateResponse struct {
	Id uuid.UUID `json:"id"`
}

type UpdateUseCase struct {
	productRepository domain.ProductRepository
}

func NewUpdateUseCase(productRepository domain.ProductRepository) *UpdateUseCase {
	return &UpdateUseCase{
		productRepository: productRepository,
	}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, r UpdateRequest) (UpdateResponse, error) {
	p, err := uc.productRepository.FetchById(ctx, r.Id)
	if err != nil {
		return UpdateResponse{}, err
	}

	p.Name = r.Name
	p.Price = r.Price
	p.UpdatedAt = time.Now()

	err = uc.productRepository.Update(ctx, p)
	if err != nil {
		return UpdateResponse{}, err
	}

	return UpdateResponse{
		Id: p.Id,
	}, nil
}
