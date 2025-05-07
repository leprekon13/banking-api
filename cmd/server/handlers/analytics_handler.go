package handlers

import (
	"banking-api/internal/db"
	"banking-api/internal/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AnalyticsHandler called")

	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "userID отсутствует", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "невалидный userID", http.StatusBadRequest)
		return
	}

	dbConn, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		http.Error(w, "подключение к БД не найдено", http.StatusInternalServerError)
		return
	}

	// вызов функции аналитики
	stats, err := db.GetMonthlyStats(dbConn, userID)
	if err != nil {
		http.Error(w, "ошибка аналитики: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
func GetKeyRateHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := services.GetCentralBankKeyRate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения ставки: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"key_rate": rate})
}
