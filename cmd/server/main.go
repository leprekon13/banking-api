package main

import (
	"banking-api/cmd/server/handlers"
	middleware "banking-api/cmd/server/middleware"
	"banking-api/config"
	"banking-api/internal/db"
	"banking-api/internal/services"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
	log.Info("logrus настроен и работает")

	cfg := config.LoadConfig()
	conn := db.ConnectDB(cfg)
	services.StartCreditScheduler(conn)

	r := mux.NewRouter()

	// Подключение к базе данных в контекст
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", conn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Публичные маршруты
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	}).Methods("GET")

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Защищённые маршруты
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/accounts", handlers.GetAccountsHandler).Methods("GET")
	protected.HandleFunc("/account", handlers.GetAccountByIDHandler).Methods("GET")
	protected.HandleFunc("/accounts", handlers.CreateAccountHandler).Methods("POST")

	protected.HandleFunc("/transfer", handlers.TransferFundsHandler).Methods("POST")

	protected.HandleFunc("/credits", handlers.CreateCreditHandler).Methods("POST")
	protected.HandleFunc("/credits/pay", handlers.PayCreditInstallmentHandler).Methods("POST")
	protected.HandleFunc("/credits/schedule", handlers.GetPaymentScheduleHandler).Methods("GET")

	protected.HandleFunc("/cards", handlers.CreateCardHandler).Methods("POST")
	protected.HandleFunc("/cards", handlers.GetCardsHandler).Methods("GET")

	protected.HandleFunc("/analytics", handlers.GetAnalyticsHandler).Methods("GET")
	protected.HandleFunc("/accounts/deposit", handlers.DepositHandler).Methods("POST")
	protected.HandleFunc("/cbr/keyrate", handlers.GetKeyRateHandler).Methods("GET")

	log.Info("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
