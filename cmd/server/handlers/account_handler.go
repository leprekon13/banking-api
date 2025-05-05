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
