package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	HabitID uint
	Habit   Habit

	Date            time.Time
	DurationMinutes uint
	Notes           string
}
