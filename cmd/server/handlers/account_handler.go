package handlers

import (
	"banking-api/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Неверный user_id", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	account, err := db.CreateAccountService(database, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания счета: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "user_id должен быть числом", http.StatusBadRequest)
		return
	}

	conn := r.Context().Value("db").(*sql.DB)
	accounts, err := db.GetAccountsByUserID(conn, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка получения счетов: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
