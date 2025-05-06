package handlers

import (
	"banking-api/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccountID int `json:"account_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Отсутствует заголовок авторизации", http.StatusUnauthorized)
		return
	}
	tokenStr := authHeader[len("Bearer "):]
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("mHbH5mvLJSfwE+YJXJtM6MwAS1vT6bf+Yp7C3Rst4aU="), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["sub"].(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный user_id в токене", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	card, err := db.CreateCardService(database, userID, input.AccountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания карты: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
func GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Отсутствует заголовок авторизации", http.StatusUnauthorized)
		return
	}
	tokenStr := authHeader[len("Bearer "):]
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("mHbH5mvLJSfwE+YJXJtM6MwAS1vT6bf+Yp7C3Rst4aU="), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["sub"].(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный user_id в токене", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	cards, err := db.GetCardsByUserID(database, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения карт: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
