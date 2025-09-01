package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	WebServerPort string
	JWTSecret     string
	JWTExpiresIn  string
}

func LoadConfig() *config {
	return &config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		WebServerPort: os.Getenv("WEB_SERVER_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiresIn:  os.Getenv("JWT_EXPIRES_IN"),
	}
}
