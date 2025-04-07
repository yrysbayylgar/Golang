package connections

import (
	"fmt"
	"myapp/internal/app/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

type Connections struct {
	DB *sqlx.DB
}

// NewConnections устанавливает соединения (например, с БД)
func NewConnections(cfg *config.Config) (*Connections, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}
	return &Connections{DB: db}, nil
}

// Close закрывает соединения
func (c *Connections) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
