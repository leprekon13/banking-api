package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
)

func CreateAccountService(db *sql.DB, userID int) (*models.Account, error) {
	account, err := CreateAccount(db, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания счета: %v", err)
	}

	return account, nil
}
