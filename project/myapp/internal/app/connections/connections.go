package connections

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver

	"myapp/internal/app/config"
)

type Connections struct {
	DB *sqlx.DB
}

func NewConnections(cfg *config.Config) (*Connections, error) {
	db, err := sqlx.Connect("postgres", cfg.DB.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return &Connections{
		DB: db,
	}, nil
}

func (c *Connections) Close() {
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			log.Println("Error closing DB:", err)
		}
	}
}
