package main

import (
	"flag"
	"myapp/internal/app"
)

func main() {
	// Парсим аргументы командной строки
	configFile := flag.String("config", "./configs/.env", "Path to configuration file")
	flag.Parse()

	// Запускаем приложение
	app.Run(*configFile)
}
