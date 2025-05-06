package app

import (
	"log"
	"myapp/internal/app/config"
	"myapp/internal/app/connections"
)

func Run(configFiles ...string) {
	// Загружаем конфиг
	cfg, err := config.NewConfig(configFiles...)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Подключаемся к БД
	conns, err := connections.NewConnections(cfg)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer conns.Close()

	log.Println("App started!")
	// Тут дальше — запуск HTTP-сервера и т.п.
}