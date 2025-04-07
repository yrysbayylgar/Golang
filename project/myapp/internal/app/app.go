package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myapp/internal/app/config"
	"myapp/internal/app/connections"
)

func Run(configFiles ...string) {
	ctx := context.Background()

	cfg, err := config.LoadConfig(configFiles...)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	conn, err := connections.NewConnections(cfg)
	if err != nil {
		log.Fatalf("Ошибка соединения: %v", err)
	}
	defer conn.Close()

	fmt.Println("Сервис запущен на порту:", cfg.Port)

	// Инициализация HTTP-сервера
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: nil, // Установите здесь свой HTTP-обработчик
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Обработка сигналов для graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	fmt.Println("Сервис завершил работу")
}
  
