package db

import (
	"database/sql"
	"fmt"
)

type Analytics struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

func GetAnalytics(db *sql.DB, userID int) (*Analytics, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount END), 0) as income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN amount END), 0) as expense
		FROM transactions
		WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
	`

	var income, expense float64
	err := db.QueryRow(query, userID).Scan(&income, &expense)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения аналитики: %v", err)
	}

	return &Analytics{
		TotalIncome:  income,
		TotalExpense: -expense, // делаем положительным
	}, nil
}

type MonthlyStats struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func GetMonthlyStats(db *sql.DB, userID int) (*MonthlyStats, error) {
	stats := &MonthlyStats{}

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN transaction_type = 'refund' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN transaction_type IN ('payment', 'transfer') THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE from_account_id IN (
			SELECT id FROM accounts WHERE user_id = $1
		)
		  AND date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
	`

	err := db.QueryRow(query, userID).Scan(&stats.Income, &stats.Expense)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
