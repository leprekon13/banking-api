package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
	"math"
	"time"
)

func CreateCreditService(db *sql.DB, userID int, amount float64, interestRate float64, months int) (*models.Credit, error) {
	startDate := time.Now()

	credit := &models.Credit{
		UserID:       userID,
		Amount:       amount,
		InterestRate: interestRate,
		StartDate:    startDate,
		Months:       months,
		CreatedAt:    time.Now(),
	}

	creditID, err := CreateCredit(db, credit)
	if err != nil {
		return nil, err
	}
	credit.ID = creditID

	monthlyRate := interestRate / 100 / 12
	payment := amount * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) / (math.Pow(1+monthlyRate, float64(months)) - 1)

	err = CreatePaymentSchedule(db, creditID, payment, months, startDate)
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func PayCreditInstallment(db *sql.DB, accountID, creditID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	var amount float64
	var scheduleID int

	err = tx.QueryRow(`
		SELECT id, amount FROM payment_schedules
		WHERE credit_id = $1 AND paid = false
		ORDER BY due_date ASC
		LIMIT 1
	`, creditID).Scan(&scheduleID, &amount)

	if err == sql.ErrNoRows {
		return fmt.Errorf("все платежи по кредиту уже оплачены")
	} else if err != nil {
		return fmt.Errorf("ошибка при получении платежа: %v", err)
	}

	var balance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, accountID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("ошибка при получении баланса: %v", err)
	}

	if balance < amount {
		return fmt.Errorf("недостаточно средств на счете")
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, accountID)
	if err != nil {
		return fmt.Errorf("ошибка при списании средств: %v", err)
	}

	_, err = tx.Exec(`UPDATE payment_schedules SET paid = true WHERE id = $1`, scheduleID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса платежа: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %v", err)
	}

	return nil
}
