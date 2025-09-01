package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/zippy/config"
	"github.com/rosset7i/zippy/internal/dto"
	"github.com/rosset7i/zippy/internal/entity"
	"github.com/rosset7i/zippy/internal/infra"
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
	if err := http.ListenAndServe(config.WebServerPort, nil); err != nil {
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

func userGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	userRepository := infra.UserRepository{
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

	userRepository := infra.UserRepository{
		DB: db,
	}

	err = userRepository.NewUser(user)
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
