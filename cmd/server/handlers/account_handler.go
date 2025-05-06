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

	w.Header().Set("Content-Type", "application/json")
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

	database := r.Context().Value("db").(*sql.DB)

	accounts, err := db.GetAccountsByUserID(database, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения счетов: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func GetAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.URL.Query().Get("account_id"))
	if err != nil {
		http.Error(w, "Неверный account_id", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	account, err := db.GetAccountByID(database, accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения счета: %v", err), http.StatusInternalServerError)
		return
	}

	if account == nil {
		http.Error(w, "Счет не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	database := r.Context().Value("db").(*sql.DB)

	var input struct {
		AccountID int     `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Amount <= 0 {
		http.Error(w, "Сумма должна быть положительной", http.StatusBadRequest)
		return
	}

	err := db.DepositToAccount(database, input.AccountID, input.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка пополнения: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Счёт успешно пополнен"))
}
