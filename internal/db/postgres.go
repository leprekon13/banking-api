package db

import (
	"database/sql"
	"fmt"
	"log"

	"banking-api/config"
	_ "github.com/lib/pq"
)

// ConnectDB открывает соединение с PostgreSQL
func ConnectDB(cfg *config.Config) *sql.DB {
	// Формируем строку подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal("База данных недоступна:", err)
	}

	log.Println("✅ Подключение к базе данных успешно")
	return db
}
