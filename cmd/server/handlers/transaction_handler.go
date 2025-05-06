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

func TransferFundsHandler(w http.ResponseWriter, r *http.Request) {
	senderID, err := strconv.Atoi(r.URL.Query().Get("from_account_id"))
	if err != nil {
		http.Error(w, "Неверный from_account_id", http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.Atoi(r.URL.Query().Get("to_account_id"))
	if err != nil {
		http.Error(w, "Неверный to_account_id", http.StatusBadRequest)
		return
	}

	var body struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	conn := r.Context().Value("db").(*sql.DB)

	err = db.TransferFunds(conn, senderID, receiverID, body.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка перевода средств: %v", err), http.StatusInternalServerError)
		return
	}

	// Получаем email отправителя
	user, err := db.GetUserByID(conn, senderID)
	if err == nil {
		// Если email есть, отправляем уведомление
		_ = services.SendPaymentEmail(user.Email, body.Amount)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Перевод выполнен"))
}
