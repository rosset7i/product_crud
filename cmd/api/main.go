package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/zippy/config"
	"github.com/rosset7i/zippy/internal/dto"
	"github.com/rosset7i/zippy/internal/entity"
	"github.com/rosset7i/zippy/internal/infra/database"
)

func main() {
	config := config.LoadConfig()

	connectionString := fmt.Sprintf(
		"dbname=%v user=%v password=%v host=%v port=%v sslmode=disable client_encoding=UTF8",
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	)

	db, err := sql.Open(config.DBDriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		userHandlers(w, r, db)
	})
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		productHandlers(w, r, db)
	})
	if err := http.ListenAndServe(config.WebServerAddress, nil); err != nil {
		log.Fatal(err)
	}
}

func userHandlers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		userGet(w, r, db)
	}
	if r.Method == http.MethodPost {
		userPost(w, r, db)
	}
}

func productHandlers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		productGetPaged(w, r, db)
	}
}

func productGetPaged(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		http.Error(w, "Invalid paging params", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	sortBy := r.URL.Query().Get("sortBy")
	if err != nil || sortBy == "" {
		http.Error(w, "Invalid paging params", http.StatusBadRequest)
		return
	}

	userRepository := database.Product{
		DB: db,
	}

	products, err := userRepository.FetchPaged(pageNumber, pageSize, sortBy)
	if err != nil {
		http.Error(w, "Products not found", http.StatusNotFound)
		return
	}

	productsDtos := make([]dto.FetchProductResponse, 0)
	for _, product := range products {
		productsDtos = append(productsDtos, dto.FetchProductResponse{
			Id:        product.Id.String(),
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		})
	}

	response, err := json.Marshal(productsDtos)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func userGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	userRepository := database.User{
		DB: db,
	}

	user, err := userRepository.FetchUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(dto.FetchUserResponse{Id: user.Id.String(), Name: user.Name, Email: user.Name})
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userRequest dto.CreateUserRequest
	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(userRequest.Name, userRequest.Email, userRequest.Password)
	if err != nil {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	userRepository := database.User{
		DB: db,
	}

	err = userRepository.Create(user)
	if err != nil {
		http.Error(w, "Error while saving the user", http.StatusBadRequest)
		return
	}

	userResponse := dto.CreateUserResponse{
		Id: user.Id.String(),
	}

	response, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
