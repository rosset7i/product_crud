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
	productRepository domain.ProductRepositoryInterface
}

func NewUpdateUseCase(productRepository domain.ProductRepositoryInterface) *UpdateUseCase {
	return &UpdateUseCase{
		productRepository: productRepository,
	}
}

func (uc *UpdateUseCase) Execute(r UpdateRequest) (UpdateResponse, error) {
	p, err := uc.productRepository.FetchById(context.TODO(), r.Id)
	if err != nil {
		return UpdateResponse{}, nil
	}

	p.Name = r.Name
	p.Price = r.Price
	p.UpdatedAt = time.Now()

	err = uc.productRepository.Update(context.TODO(), p)
	if err != nil {
		return UpdateResponse{}, nil
	}

	return UpdateResponse{
		Id: p.Id,
	}, nil
}
