package models

import (
	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model

	Name        string
	Description string
	Frequency   string

	Sessions []Session
}
