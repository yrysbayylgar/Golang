package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        string `env:"PORT" envDefault:"8080"`
}

// LoadConfig загружает конфигурацию из .env и переменных окружения
func LoadConfig(filenames ...string) (*Config, error) {
	_ = godotenv.Load(filenames...) // Загружаем .env (если есть)
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	return &cfg, nil
}
