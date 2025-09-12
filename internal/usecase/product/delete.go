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
	productRepository domain.ProductRepository
}

func NewDeleteUseCase(productRepository domain.ProductRepository) *DeleteUseCase {
	return &DeleteUseCase{
		productRepository: productRepository,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, r DeleteRequest) (DeleteResponse, error) {
	err := uc.productRepository.Delete(ctx, r.Id)
	if err != nil {
		return DeleteResponse{}, err
	}

	return DeleteResponse(r), nil
}
