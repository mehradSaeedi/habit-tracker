package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	HabitID uint  `json:"habitId"`
	Habit   Habit `json:"habit,omitzero"`

	Date            time.Time `json:"date"`
	DurationMinutes uint      `json:"durationMinutes"`
	Notes           string    `json:"notes"`
}
