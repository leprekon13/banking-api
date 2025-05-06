package services

import (
	"banking-api/internal/db"
	"database/sql"
	"log"
	"time"
)

// StartCreditScheduler запускает фоновую задачу раз в 12 часов
func StartCreditScheduler(database *sql.DB) {
	go func() {
		for {
			log.Println("Шедулер: запуск проверки просроченных платежей...")

			err := db.ProcessOverduePayments(database)
			if err != nil {
				log.Printf("Шедулер: ошибка при обработке просроченных платежей: %v", err)
			} else {
				log.Println("Шедулер: проверка завершена успешно")
			}

			time.Sleep(12 * time.Hour)
		}
	}()
}
