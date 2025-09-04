package config

import (
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	_ "github.com/joho/godotenv/autoload"
)

var cfg *Config

type Config struct {
	DBDriver         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	WebServerAddress string
	JWTSecret        string
	JWTExpiresIn     int
	TokenAuth        *jwtauth.JWTAuth
}

func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}

	jwtExpiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		panic("Could not load env vars")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	cfg = &Config{
		DBDriver:         os.Getenv("DB_DRIVER"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		WebServerAddress: os.Getenv("WEB_SERVER_ADDRESS"),
		JWTSecret:        jwtSecret,
		JWTExpiresIn:     jwtExpiresIn,
		TokenAuth:        tokenAuth,
	}

	return cfg
}
