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

type HabitStats struct {
	TotalSessions  int     `json:"total_sessions"`
	TotalMinutes   int     `json:"total_minutes"`
	AverageMinutes float64 `json:"average_minutes"`
	LongestSession int     `json:"longest_session"`
}
