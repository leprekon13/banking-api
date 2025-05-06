package models

import "time"

type Card struct {
	ID             int       `json:"id"`
	AccountID      int       `json:"account_id"`
	CardNumber     string    `json:"card_number"`
	ExpirationDate time.Time `json:"expiration_date"`
	CVV            string    `json:"cvv"`
	HMAC           string    `json:"hmac"`
	CreatedAt      time.Time `json:"created_at"`
}
