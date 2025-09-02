package database

import (
	"database/sql"

	"github.com/rosset7i/zippy/internal/entity"
)

type User struct {
	DB *sql.DB
}

func (u *User) Create(user *entity.User) error {
	_, err := u.DB.Exec(
		"INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)",
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (u *User) FetchUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	rows := u.DB.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1", email)
	err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return &user, err
}
