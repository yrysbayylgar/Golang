package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию сервиса
type Config struct {
	Port        int    // Порт сервиса
	DatabaseURL string // URL базы данных
	JWTSecret   string // Секрет для JWT
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загрузка порта
	port, err := strconv.Atoi(getEnv("PORT", "8081"))
	if err != nil {
		return nil, errors.New("недопустимое значение PORT")
	}

	// Загрузка URL базы данных - с поддержкой отдельных параметров подключения
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		// Если DATABASE_URL не задан, пытаемся составить его из отдельных компонентов
		dbHost := getEnv("DB_HOST", "")
		dbPort := getEnv("DB_PORT", "")
		dbUser := getEnv("DB_USER", "")
		dbPass := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "")

		// Проверяем, что все необходимые компоненты предоставлены
		if dbHost != "" && dbPort != "" && dbUser != "" && dbName != "" {
			// Конструируем строку подключения PostgreSQL
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				dbUser, dbPass, dbHost, dbPort, dbName)
		} else {
			return nil, errors.New("требуется DATABASE_URL или все DB_* параметры")
		}
	}

	// Загрузка JWT секрета
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("требуется JWT_SECRET")
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}, nil
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
