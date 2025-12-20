package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Вспомогательная функция для переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Connect() *pgxpool.Pool {
	ctx := context.Background()

	// Получаем настройки из переменных окружения (для Docker)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "Sports_supplement_store")

	// Формируем DSN
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err, "Ошибка парсинга конфигурации базы данных")
	}
	config.MaxConns = 5
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err, "Ошибка подключения к базе данных")
	}

	// Проверка подключения
	ctx_for_check, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx_for_check); err != nil {
		log.Fatal(err, "Ошибка пинга базы данных")
	}
	fmt.Println("Успешное подключение к базе данных")

	return pool
}
