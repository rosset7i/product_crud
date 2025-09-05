package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	_ "github.com/joho/godotenv/autoload"
)

var cfg *Config

type Config struct {
	DBDriver         string
	ConnectionString string
	WebServerAddress string
	JWTSecret        string
	JWTExpiresIn     int
	TokenAuth        *jwtauth.JWTAuth
}

var (
	ErrCouldNotLoadEnvVars = errors.New("could not load all required environment variables")
)

func LoadConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	dBDriver := os.Getenv("DB_DRIVER")
	dBHost := os.Getenv("DB_HOST")
	dBPort := os.Getenv("DB_PORT")
	dBUser := os.Getenv("DB_USER")
	dBPassword := os.Getenv("DB_PASSWORD")
	dBName := os.Getenv("DB_NAME")
	webServerAddress := os.Getenv("WEB_SERVER_ADDRESS")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresInEnv := os.Getenv("JWT_EXPIRES_IN")

	jwtExpiresIn, err := strconv.Atoi(jwtExpiresInEnv)
	if err != nil || jwtSecret == "" {
		return nil, fmt.Errorf("%w: invalid JWT config", ErrCouldNotLoadEnvVars)
	}

	if dBDriver == "" || dBHost == "" || dBPort == "" || dBUser == "" || dBPassword == "" || dBName == "" {
		return nil, fmt.Errorf("%w: missing database configuration", ErrCouldNotLoadEnvVars)
	}

	connectionString := fmt.Sprintf(
		"dbname=%v user=%v password=%v host=%v port=%v sslmode=disable client_encoding=UTF8",
		dBName,
		dBUser,
		dBPassword,
		dBHost,
		dBPort,
	)

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	cfg = &Config{
		DBDriver:         dBDriver,
		ConnectionString: connectionString,
		WebServerAddress: webServerAddress,
		JWTSecret:        jwtSecret,
		JWTExpiresIn:     jwtExpiresIn,
		TokenAuth:        tokenAuth,
	}

	return cfg, nil
}
