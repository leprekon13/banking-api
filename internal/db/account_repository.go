package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateAccount(db *sql.DB, userID int) (*models.Account, error) {
	account := &models.Account{
		UserID:    userID,
		Balance:   0.0, // Начальный баланс
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := db.QueryRow(`
        INSERT INTO accounts (user_id, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4) RETURNING id`, account.UserID, account.Balance, account.CreatedAt, account.UpdatedAt).Scan(&account.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании счета: %v", err)
	}

	return account, nil
}
