package config

import (
	"log"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer HTTPServerConfig
	DB         DBConfig
}

type HTTPServerConfig struct {
	Host string `env:"HTTP_HOST" envDefault:"localhost"`
	Port string `env:"HTTP_PORT" envDefault:"8080"`
}

type DBConfig struct {
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

// NewConfig загружает переменные окружения и возвращает конфигурацию
func NewConfig(filenames ...string) (*Config, error) {
	for _, f := range filenames {
		if err := godotenv.Load(f); err != nil {
			log.Printf("No .env file found at %s", f)
		}
	}
	
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}