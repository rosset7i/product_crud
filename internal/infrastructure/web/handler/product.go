package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rosset7i/product_crud/internal/infrastructure/web"
	"github.com/rosset7i/product_crud/internal/usecase/product"
)

type ProductHandler struct {
	fetchPagedProductsUseCase *product.FetchPagedProductsUseCase
	fetchByIdUseCase          *product.FetchByIdUseCase
	createUseCase             *product.CreateUseCase
	updateUseCase             *product.UpdateUseCase
	deleteUseCase             *product.DeleteUseCase
}

func NewProductHandler(
	fetchPagedProductsUseCase *product.FetchPagedProductsUseCase,
	fetchByIdUseCase *product.FetchByIdUseCase,
	createUseCase *product.CreateUseCase,
	updateUseCase *product.UpdateUseCase,
	deleteUseCase *product.DeleteUseCase,
) *ProductHandler {
	return &ProductHandler{
		fetchPagedProductsUseCase: fetchPagedProductsUseCase,
		fetchByIdUseCase:          fetchByIdUseCase,
		createUseCase:             createUseCase,
		updateUseCase:             updateUseCase,
		deleteUseCase:             deleteUseCase,
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
// @Failure      400         {object}  web.errorResponse
// @Failure      422         {object}  web.errorResponse
// @Router       /v1/products [get]
// @Security Bearer
func (h *ProductHandler) FetchPaged(w http.ResponseWriter, r *http.Request) {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	sort := r.URL.Query().Get("sort")
	if err != nil || sort == "" {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.fetchPagedProductsUseCase.Execute(r.Context(), product.FetchPagedProductsRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Sort:       sort,
	})
	if err != nil {
		web.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusOK, response)
}

// GetProduct godoc
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true "id"
// @Success      200  {object}  product.FetchByIdResponse
// @Failure      400  {object}  web.errorResponse
// @Failure      404  {object}  web.errorResponse
// @Router       /v1/products/{id} [get]
// @Security Bearer
func (h *ProductHandler) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.fetchByIdUseCase.Execute(r.Context(), product.FetchByIdRequest{Id: id})
	if err != nil {
		web.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusOK, response)
}

// Create Product godoc
// @Tags         products
// @Param        request  body      product.CreateRequest  true "payload"
// @Success      201      {object}  product.CreateResponse
// @Failure      400      {object}  web.errorResponse
// @Failure      422      {object}  web.errorResponse
// @Router       /v1/products [post]
// @Security Bearer
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeJSONBody[product.CreateRequest](r)
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.createUseCase.Execute(r.Context(), req)
	if err != nil {
		web.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusCreated, response)
}

// UpdateProduct godoc
// @Tags         products
// @Param        request  body      product.UpdateRequest  true "payload"
// @Success      200      {object}  product.UpdateResponse
// @Failure      400      {object}  web.errorResponse
// @Failure      422      {object}  web.errorResponse
// @Router       /v1/products [put]
// @Security Bearer
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeJSONBody[product.UpdateRequest](r)
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.updateUseCase.Execute(r.Context(), req)
	if err != nil {
		web.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusOK, response)
}

// DeleteProduct godoc
// @Tags         products
// @Param        id   query     string  true "id"
// @Success      200  {object}  product.DeleteResponse
// @Failure      400  {object}  web.errorResponse
// @Failure      422  {object}  web.errorResponse
// @Router       /v1/products [delete]
// @Security Bearer
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		web.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.deleteUseCase.Execute(r.Context(), product.DeleteRequest{Id: id})
	if err != nil {
		web.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	web.WriteJSON(w, http.StatusOK, response)
}
