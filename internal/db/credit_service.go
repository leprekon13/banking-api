package db

import (
	"banking-api/models"
	"database/sql"
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
