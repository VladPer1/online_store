package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"online_store/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunServerWithShutdown(db *pgxpool.Pool, mux *http.ServeMux) {
	// Создаем свой мукс

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Используем наш мукс
	}

	// Регистрация обработчиков
	handlers.RegisterRoutes(mux, db)

	// Запуск сервера в горутине
	go func() {
		log.Println("Сервер запущен на http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Graceful shutdown с таймаутом 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка graceful shutdown: %v", err)
	} else {
		log.Println("Сервер корректно остановлен")
	}
}
