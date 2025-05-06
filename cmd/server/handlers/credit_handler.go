package handlers

import (
	"banking-api/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateCreditHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID       int     `json:"user_id"`
		Amount       float64 `json:"amount"`
		InterestRate float64 `json:"interest_rate"`
		Months       int     `json:"months"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	credit, err := db.CreateCreditService(database, input.UserID, input.Amount, input.InterestRate, input.Months)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания кредита: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}
