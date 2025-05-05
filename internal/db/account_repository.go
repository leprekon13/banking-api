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

func GetAccountsByUserID(db *sql.DB, userID int) ([]models.Account, error) {
	rows, err := db.Query("SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса счетов: %v", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func GetAccountByID(db *sql.DB, accountID int) (*models.Account, error) {
	row := db.QueryRow("SELECT id, user_id, balance, created_at FROM accounts WHERE id = $1", accountID)

	var acc models.Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найден
		}
		return nil, err
	}

	return &acc, nil
}
