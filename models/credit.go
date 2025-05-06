package models

import "time"

type Credit struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	StartDate    time.Time `json:"start_date"`
	Months       int       `json:"months"`
	CreatedAt    time.Time `json:"created_at"`
}

type PaymentSchedule struct {
	ID        int       `json:"id"`
	CreditID  int       `json:"credit_id"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Paid      bool      `json:"paid"`
	CreatedAt time.Time `json:"created_at"`
}
