package config

import (
	"log"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
)

type Conf struct {
	Auth   ConfAuth
	Server ConfServer
	DB     ConfDB
}

type ConfAuth struct {
	JwtSecret    string        `env:"JWT_SECRET,required"`
	JwtExpiresIn time.Duration `env:"JWT_EXPIRES_IN,required"`
	JwtAuth      *jwtauth.JWTAuth
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	Debug    bool   `env:"DB_DEBUG,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	c.Auth.JwtAuth = jwtauth.New("HS256", []byte(c.Auth.JwtSecret), nil)

	return &c
}
