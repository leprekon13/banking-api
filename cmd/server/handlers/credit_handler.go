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

func CreateCreditHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccountID    int     `json:"account_id"`
		Amount       float64 `json:"amount"`
		InterestRate float64 `json:"interest_rate"`
		Months       int     `json:"duration_months"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
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

	credit, err := db.CreateCreditService(database, userID, input.Amount, input.InterestRate, input.Months)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания кредита: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func PayCreditInstallmentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CreditID int `json:"credit_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	err := db.PayNextInstallment(database, input.CreditID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка оплаты платежа: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ближайший платеж успешно оплачен"))
}

func GetPaymentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	creditIDStr := r.URL.Query().Get("credit_id")
	if creditIDStr == "" {
		http.Error(w, "credit_id обязателен", http.StatusBadRequest)
		return
	}

	creditID, err := strconv.Atoi(creditIDStr)
	if err != nil {
		http.Error(w, "credit_id должен быть числом", http.StatusBadRequest)
		return
	}

	database := r.Context().Value("db").(*sql.DB)

	schedule, err := db.GetPaymentScheduleByCreditID(database, creditID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения графика: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
