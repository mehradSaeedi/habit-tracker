package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mehradSaeedi/habit-tracker/pkg/config"
	"github.com/mehradSaeedi/habit-tracker/pkg/models"
)

// ---CREATE---
func CreateHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	err := json.NewDecoder(r.Body).Decode(&habit)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := config.DB.Create(&habit)
	if result.Error != nil {
		http.Error(w, "Could not create habit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)

}

// ---READ---
func GetHabits(w http.ResponseWriter, r *http.Request) {
	var habits []models.Habit

	result := config.DB.Find(&habits)
	if result.Error != nil {
		http.Error(w, "Could not fetch habits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

func GetHabitByID(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result := config.DB.Preload("Sessions").First(&habit, id)
	if result.Error != nil {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
}

// ---UPDATE---
func UpdateHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	var updateHabit models.Habit
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result := config.DB.First(&habit, id)
	if result.RowsAffected == 0 {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	// We use a temporary struct to decode into because the incoming JSON might overwrite the fields that
	// we don't want touched (like IDs and timestapms)
	err = json.NewDecoder(r.Body).Decode(&updateHabit)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	// Copy only the fileds we allow the client to update
	habit.Name = updateHabit.Name
	habit.Description = updateHabit.Description
	habit.Frequency = updateHabit.Frequency

	result = config.DB.Save(&habit)
	if result.Error != nil {
		http.Error(w, "Could not update habit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
}

// ---DELETE---
func DeleteHabit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Habit{}, id)

	if result.Error != nil {
		http.Error(w, "Could not delete habit", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "habit deleted",
	})

}

// ---STATS---
func GetHabitStats(w http.ResponseWriter, r *http.Request) {
	habitIdStr := r.PathValue("id")
	habitId, err := strconv.Atoi(habitIdStr)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	var stats models.HabitStats

	result := config.DB.Model(&models.Session{}).Select(`
	COUNT(*) AS total_sessions,
	COALESCE(SUM(duration_minutes), 0) AS total_minutes,
	COALESCE(AVG(duration_minutes), 0) AS average_minutes,
	COALESCE(MAX(duration_minutes), 0) AS longest_session
	`).Where("habit_id = ?", habitId).Scan(&stats)
	if result.Error != nil {
		http.Error(w, "Could not calculate statistics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Could not encode statistics JSON", http.StatusInternalServerError)
		return
	}
}

func ClearSessions(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	result := config.DB.First(&models.Habit{}, id)
	if result.RowsAffected == 0 {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	result = config.DB.Where("habit_id = ?", id).Delete(&models.Session{})
	if result.Error != nil {
		http.Error(w, "Could not delete sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sessions Cleared",
	})
}
