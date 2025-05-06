package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"myapp/internal/config"
	"myapp/internal/handlers"
	"myapp/internal/middleware"
	"myapp/internal/repository/postgres"
	"myapp/internal/service"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Подключение к базе данных
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Настройка репозитория, сервиса и обработчика
	repo := postgres.NewRepository(db)
	svc := service.NewService(repo)
	handler := handlers.NewHandler(svc)

	// Настройка маршрутизатора
	router := mux.NewRouter()
	
	// Применение промежуточного ПО
	router.Use(middleware.Logger)
	router.Use(middleware.JSONContentType)
	
	// Промежуточное ПО аутентификации для защищенных маршрутов
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(middleware.JWTAuth(cfg.JWTSecret))

	// Регистрация маршрутов
	handler.RegisterRoutes(authRouter)
	
	// Маршрут проверки работоспособности
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Group Service работает нормально"))
	}).Methods("GET")

	// Создание сервера
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Запуск Group Service на порту %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось запустить сервер: %v", err)
		}
	}()

	// Ожидание сигнала прерывания
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Создание крайнего срока для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Завершение работы сервера
	log.Println("Завершение работы Group Service...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение работы сервера: %v", err)
	}

	log.Println("Group Service остановлен корректно")
}