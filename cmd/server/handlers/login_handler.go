package handlers

import (
	"banking-api/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)
	token, err := db.LoginUserService(database, req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка входа: %v", err), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
