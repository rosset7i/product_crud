package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var cfg *config

type config struct {
	DBDriver         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	WebServerAddress string
	JWTSecret        string
	JWTExpiresIn     string
}

func LoadConfig() *config {
	if cfg != nil {
		return cfg
	}

	cfg = &config{
		DBDriver:         os.Getenv("DB_DRIVER"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		WebServerAddress: os.Getenv("WEB_SERVER_ADDRESS"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpiresIn:     os.Getenv("JWT_EXPIRES_IN"),
	}

	return cfg
}
