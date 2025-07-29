package domain

import (
	"gorm.io/gorm"
	"time"
)

type Chronicle struct {
	gorm.Model
	Date     time.Time
	Event    string
	Priority int
}
