package web

import (
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
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        pageNumber  query     int     true "pageNumber"
// @Param        pageSize    query     int     true "pageSize"
// @Param        sort        query     string  true "sort"
// @Success      200         {object}  product.FetchPagedProductsResponse
// @Failure      400         {object}  errorResponse
// @Failure      422         {object}  errorResponse
// @Router       /v1/products [get]
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

	response, err := product.NewFetchPagedProductsUseCase(h.productRepository).Execute(r.Context(), product.FetchPagedProductsRequest{
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
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true "id"
// @Success      200  {object}  product.FetchByIdResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Router       /v1/products/{id} [get]
// @Security Bearer
func (h *ProductHandler) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewFetchByIdUseCase(h.productRepository).Execute(r.Context(), product.FetchByIdRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// Create Product godoc
// @Tags         products
// @Param        request  body      product.CreateRequest  true "payload"
// @Success      201      {object}  product.CreateResponse
// @Failure      400      {object}  errorResponse
// @Failure      422      {object}  errorResponse
// @Router       /v1/products [post]
// @Security Bearer
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSONBody[product.CreateRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewCreateUseCase(h.productRepository).Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

// UpdateProduct godoc
// @Tags         products
// @Param        request  body      product.UpdateRequest  true "payload"
// @Success      200      {object}  product.UpdateResponse
// @Failure      400      {object}  errorResponse
// @Failure      422      {object}  errorResponse
// @Router       /v1/products [put]
// @Security Bearer
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSONBody[product.UpdateRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewUpdateUseCase(h.productRepository).Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// DeleteProduct godoc
// @Tags         products
// @Param        id   query     string  true "id"
// @Success      200  {object}  product.DeleteResponse
// @Failure      400  {object}  errorResponse
// @Failure      422  {object}  errorResponse
// @Router       /v1/products [delete]
// @Security Bearer
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := product.NewDeleteUseCase(h.productRepository).Execute(r.Context(), product.DeleteRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}
