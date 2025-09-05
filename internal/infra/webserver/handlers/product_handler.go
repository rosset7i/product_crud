package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rosset7i/zippy/internal/dto"
	"github.com/rosset7i/zippy/internal/entity"
	"github.com/rosset7i/zippy/internal/infra/database"
	"github.com/rosset7i/zippy/internal/infra/webserver"
)

type ProductHandler struct {
	ProductDB *database.ProductRepository
}

func NewProductHandler(productDb *database.ProductRepository) *ProductHandler {
	return &ProductHandler{
		ProductDB: productDb,
	}
}

// List Products godoc
// @Summary      List products
// @Description  Retrieve a paginated list of products with optional sorting
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        pageNumber  query     int     false  "Page number (default: 1)"
// @Param        pageSize    query     int     false  "Number of items per page (default: 10)"
// @Param        sort        query     string  false  "Sort order: asc or desc (default: asc)"
// @Success      200         {array}   entity.Product
// @Failure      400         {object}  webserver.ErrorResponse  "Invalid query parameters"
// @Failure      500         {object}  webserver.ErrorResponse  "Internal server error"
// @Router       /products [get]
// @Security Bearer
func (h *ProductHandler) FetchPaged(w http.ResponseWriter, r *http.Request) {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	sort := r.URL.Query().Get("sort")
	if err != nil || sort == "" {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.ProductDB.FetchPaged(pageNumber, pageSize, sort)
	if err != nil {
		webserver.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusOK, products)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieve details of a single product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID (UUID)"
// @Success      200  {object}  entity.Product
// @Failure      400  {object}  webserver.ErrorResponse  "Invalid product ID"
// @Failure      404  {object}  webserver.ErrorResponse  "Product not found"
// @Failure      500  {object}  webserver.ErrorResponse  "Internal server error"
// @Router       /products/{id} [get]
// @Security Bearer
func (h *ProductHandler) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.ProductDB.FetchById(id)
	if err != nil {
		webserver.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusOK, product)
}

// Create Product godoc
// @Summary      Create a new product
// @Description  Create a new product with a name and price
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateProductRequest  true  "Product creation request"
// @Success      201      {object}  entity.Product
// @Failure      400      {object}  webserver.ErrorResponse  "Invalid request payload"
// @Failure      500      {object}  webserver.ErrorResponse  "Internal server error"
// @Router       /products [post]
// @Security Bearer
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := entity.NewProduct(request.Name, request.Price)
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		webserver.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Update product details (name and price)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateProductRequest  true  "Product update request"
// @Success      200      {object}  entity.Product
// @Failure      400      {object}  webserver.ErrorResponse  "Invalid request payload"
// @Failure      404      {object}  webserver.ErrorResponse  "Product not found"
// @Failure      500      {object}  webserver.ErrorResponse  "Internal server error"
// @Router       /products [put]
// @Security Bearer
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateProductRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.ProductDB.FetchById(request.Id)
	if err != nil {
		webserver.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	product.Name = request.Name
	product.Price = request.Price
	product.UpdatedAt = time.Now()

	err = h.ProductDB.Update(product)
	if err != nil {
		webserver.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Permanently delete a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   query     string  true  "Product ID (UUID)"
// @Success      200  {string}  string  "Product deleted successfully"
// @Failure      400  {object}  webserver.ErrorResponse   "Invalid product ID"
// @Failure      404  {object}  webserver.ErrorResponse   "Product not found"
// @Failure      500  {object}  webserver.ErrorResponse   "Internal server error"
// @Router       /products [delete]
// @Security Bearer
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		webserver.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		webserver.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	webserver.WriteJSON(w, http.StatusOK, struct {
		Id uuid.UUID `json:"id"`
	}{Id: id})
}
