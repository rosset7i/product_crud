package server

import (
	"github.com/rosset7i/product_crud/internal/infrastructure/database"
	"github.com/rosset7i/product_crud/internal/infrastructure/web/handler"
	"github.com/rosset7i/product_crud/internal/usecase/product"
	"github.com/rosset7i/product_crud/internal/usecase/user"
)

type Container struct {
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
}

func (s *Server) init() {
	// repositories
	userRepository := database.NewUserRepository(s.db)
	productRepository := database.NewProductRepository(s.db)

	// use cases
	registerUseCase := user.NewRegisterUseCase(userRepository)
	loginUseCase := user.NewLoginUseCase(userRepository, s.c.Auth.JwtAuth, s.c.Auth.JwtExpiresIn)
	fetchPagedProductsUseCase := product.NewFetchPagedProductsUseCase(productRepository)
	fetchByIdUseCase := product.NewFetchByIdUseCase(productRepository)
	createUseCase := product.NewCreateUseCase(productRepository)
	updateUseCase := product.NewUpdateUseCase(productRepository)
	deleteUseCase := product.NewDeleteUseCase(productRepository)

	// handlers
	userHandler := handler.NewUserHandler(registerUseCase, loginUseCase)
	productHandler := handler.NewProductHandler(fetchPagedProductsUseCase, fetchByIdUseCase, createUseCase, updateUseCase, deleteUseCase)

	s.container = &Container{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
}
