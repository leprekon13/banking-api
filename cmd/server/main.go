package main

import (
	"banking-api/config"
	"banking-api/internal/db"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфиг
	cfg := config.LoadConfig()

	// Подключаемся к БД
	_ = db.ConnectDB(cfg)

	// Роутер
	r := mux.NewRouter()

	// Тестовый маршрут
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
