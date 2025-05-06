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

func TransferFunds(db *sql.DB, senderID, receiverID int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", senderID).Scan(&senderBalance)
	if err != nil {
		return fmt.Errorf("не удалось получить баланс отправителя: %v", err)
	}
	if senderBalance < amount {
		return fmt.Errorf("недостаточно средств")
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, senderID)
	if err != nil {
		return fmt.Errorf("не удалось списать средства у отправителя: %v", err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		return fmt.Errorf("не удалось зачислить средства получателю: %v", err)
	}

	_, err = tx.Exec(
		"INSERT INTO transactions (from_account_id, to_account_id, amount, created_at) VALUES ($1, $2, $3, $4)",
		senderID, receiverID, amount, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("не удалось записать транзакцию: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %v", err)
	}

	return nil
}
func DepositToAccount(db *sql.DB, accountID int, amount float64) error {
	_, err := db.Exec(`
		UPDATE accounts SET balance = balance + $1 WHERE id = $2
	`, amount, accountID)
	return err
}
