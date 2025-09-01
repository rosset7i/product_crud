package infra

import (
	"database/sql"

	"github.com/rosset7i/zippy/internal/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func (u *UserRepository) NewUser(user *entity.User) error {
	_, err := u.DB.Exec(
		"INSERT INTO users (id, name, email, password_hash) VALUES ($1,$2,$3,$4)",
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (u *UserRepository) FetchUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	rows := u.DB.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email)
	err := rows.Scan(&user.Id, &user.Name, &user.Email)

	return &user, err
}
