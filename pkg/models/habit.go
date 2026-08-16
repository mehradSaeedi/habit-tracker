package models

import (
	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model

	Name        string `josn:"name"`
	Description string `josn:"description"`
	Frequency   string `josn:"frequency"`

	Sessions []Session `json:"sessions,omitempty"`
}
