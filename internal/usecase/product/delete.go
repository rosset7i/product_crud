package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/domain"
)

type DeleteRequest struct {
	Id uuid.UUID `json:"id"`
}

type DeleteResponse struct {
	Id uuid.UUID `json:"id"`
}

type DeleteUseCase struct {
	productRepository domain.ProductRepositoryInterface
}

func NewDeleteUseCase(productRepository domain.ProductRepositoryInterface) *DeleteUseCase {
	return &DeleteUseCase{
		productRepository: productRepository,
	}
}

func (uc *DeleteUseCase) Execute(r DeleteRequest) (DeleteResponse, error) {
	err := uc.productRepository.Delete(context.TODO(), r.Id)
	if err != nil {
		return DeleteResponse{}, nil
	}

	return DeleteResponse(r), nil
}
