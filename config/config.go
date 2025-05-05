package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит параметры окружения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// LoadConfig загружает переменные из .env
func LoadConfig() *Config {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: ", err)
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
