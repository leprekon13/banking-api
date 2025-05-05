package db

import (
	"banking-api/models"
	"database/sql"
	"fmt"
)

func AddUser(db *sql.DB, user *models.User) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil {
		return fmt.Errorf("ошибка проверки уникальности email: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("email уже используется")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&count)
	if err != nil {
		return fmt.Errorf("ошибка проверки уникальности username: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("имя пользователя уже занято")
	}

	_, err = db.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		user.Username, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя: %v", err)
	}

	return nil
}

func LoginUser(db *sql.DB, email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с таким email не найден")
		}
		return nil, fmt.Errorf("ошибка при поиске пользователя: %v", err)
	}
	return &user, nil
}
