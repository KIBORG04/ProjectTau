package domain

import (
	"encoding/json"
	"time"
)

// ForecastHistory хранит снимок прогноза, сделанный в определенное время.
type ForecastHistory struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time // Дата, когда прогноз был СДЕЛАН. Это наш ключ к истории.

	ForecastType string `gorm:"index"` // 'weekly' или 'daily'

	// JSON с данными прогноза: {"2025-10-20": 55.5, "2025-10-21": 56.0, ...}
	Data json.RawMessage `gorm:"type:jsonb"`
}