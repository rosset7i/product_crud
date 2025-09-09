package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/infrastructure/database"
	"github.com/rosset7i/product_crud/internal/usecase/product"
)

type ProductHandler struct {
	productRepository *database.ProductRepository
}

func NewProductHandler(productRepository *database.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepository: productRepository,
	}
}

// List Products godoc
// @Summary      List products
// @Description  Retrieve a paginated list of products with optional sorting
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        pageNumber  query     int     false "pageNumber"
// @Param        pageSize    query     int     false "pageSize"
// @Param        sort        query     string  false "sort"
// @Success      200         {object}  product.FetchPagedProductsResponse
// @Failure      400         {object}  errorResponse
// @Failure      422         {object}  errorResponse
// @Router       /products [get]
// @Security Bearer
func (h *ProductHandler) FetchPaged(w http.ResponseWriter, r *http.Request) {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	sort := r.URL.Query().Get("sort")
	if err != nil || sort == "" {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewFetchPagedProductsUseCase(h.productRepository).Execute(product.FetchPagedProductsRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Sort:       sort,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieve details of a single product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true "id"
// @Success      200  {object}  product.FetchByIdResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Router       /products/{id} [get]
// @Security Bearer
func (h *ProductHandler) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewFetchByIdUseCase(h.productRepository).Execute(product.FetchByIdRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// Create Product godoc
// @Summary      Create a new product
// @Description  Create a new product with a name and price
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      product.CreateRequest  true "CreateRequest"
// @Success      201      {object}  product.CreateResponse
// @Failure      400      {object}  errorResponse
// @Failure      422      {object}  errorResponse
// @Router       /products [post]
// @Security Bearer
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewCreateUseCase(h.productRepository).Execute(req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Update product details (name and price)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      product.UpdateRequest  true "UpdateRequest"
// @Success      200      {object}  product.UpdateResponse
// @Failure      400      {object}  errorResponse
// @Failure      422      {object}  errorResponse
// @Router       /products [put]
// @Security Bearer
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req product.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewUpdateUseCase(h.productRepository).Execute(req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Permanently delete a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   query     string  true "id"
// @Success      200  {object}  product.DeleteResponse
// @Failure      400  {object}  errorResponse
// @Failure      422  {object}  errorResponse
// @Router       /products [delete]
// @Security Bearer
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewDeleteUseCase(h.productRepository).Execute(product.DeleteRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}
