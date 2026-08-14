package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mehradSaeedi/habit-tracker/pkg/config"
	"github.com/mehradSaeedi/habit-tracker/pkg/models"
)

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

func GetHabitByID(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result := config.DB.First(&habit, id)
	if result.Error != nil {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
}

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
