package repository

import (
	"ssstatistics/internal/domain"
	"time"
)

func AddNewChronicle(dateStr, event string) error {
	// Парсим дату из строки (ожидаем формат "YYYY-MM-DD")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	newChronicle := domain.Chronicle{
		Date:  date,
		Event: event,
	}

	result := Database.Create(&newChronicle)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
