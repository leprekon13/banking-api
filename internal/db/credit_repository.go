package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateCredit(db *sql.DB, credit *models.Credit) (int, error) {
	var creditID int
	err := db.QueryRow(`
		INSERT INTO credits (user_id, amount, interest_rate, start_date, months, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, credit.UserID, credit.Amount, credit.InterestRate, credit.StartDate, credit.Months, credit.CreatedAt).Scan(&creditID)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании кредита: %v", err)
	}

	return creditID, nil
}

func CreatePaymentSchedule(db *sql.DB, creditID int, monthlyAmount float64, months int, startDate time.Time) error {
	for i := 0; i < months; i++ {
		dueDate := startDate.AddDate(0, i+1, 0)
		_, err := db.Exec(`
			INSERT INTO payment_schedules (credit_id, amount, due_date, paid, created_at)
			VALUES ($1, $2, $3, false, $4)
		`, creditID, monthlyAmount, dueDate, time.Now())

		if err != nil {
			return fmt.Errorf("ошибка при создании графика платежей: %v", err)
		}
	}
	return nil
}
