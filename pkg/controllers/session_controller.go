package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mehradSaeedi/habit-tracker/pkg/config"
	"github.com/mehradSaeedi/habit-tracker/pkg/models"
)

func CreateSession(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	result := config.DB.First(&models.Habit{}, id)
	if result.Error != nil {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	var session models.Session
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	session.HabitID = uint(id)

	result = config.DB.Create(&session)
	if result.Error != nil {
		http.Error(w, "Couldn't create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}

}

func GetSessionsForHabit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	var habit models.Habit
	if err := config.DB.First(&habit, id).Error; err != nil {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	var sessions []models.Session
	dbQuery := config.DB.Where("habit_id = ?", id)

	layout := "2006-01-02"

	var fromDate time.Time
	from := query.Get("from")
	if from != "" {
		fromDate, err = time.Parse(layout, from)
		if err != nil {
			http.Error(w, "Invalid query", http.StatusBadRequest)
			return
		}
		dbQuery = dbQuery.Where("date >= ?", fromDate)
	}
	to := query.Get("to")
	var toDate time.Time
	if to != "" {
		toDate, err = time.Parse(layout, to)
		if err != nil {
			http.Error(w, "Invalid query", http.StatusBadRequest)
			return
		}
		dbQuery = dbQuery.Where("date <= ?", toDate)
	}

	result := dbQuery.Find(&sessions)
	if result.Error != nil {
		http.Error(w, "Could not fetch sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sessions)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}

}

func GetSessionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var session models.Session
	result := config.DB.First(&session, id)
	if result.RowsAffected == 0 {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		http.Error(w, "Could not fetch session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}

}

func UpdateSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("id")
	id, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var session models.Session

	result := config.DB.First(&session, id)
	if result.Error != nil {
		http.Error(w, "Could not update session", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	var updateSession models.Session
	err = json.NewDecoder(r.Body).Decode(&updateSession)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	session.Notes = updateSession.Notes
	session.DurationMinutes = updateSession.DurationMinutes

	result = config.DB.Save(&session)
	if result.Error != nil {
		http.Error(w, "Could not update session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("id")
	id, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Session{}, id)
	if result.Error != nil {
		http.Error(w, "Could not delete session", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Could not find session", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session Deleted",
	})
}
