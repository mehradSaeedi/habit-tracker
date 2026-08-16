package routes

import (
	"net/http"

	"github.com/mehradSaeedi/habit-tracker/pkg/controllers"
)

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /habits", controllers.CreateHabit)
	mux.HandleFunc("GET /habits", controllers.GetHabits)
	mux.HandleFunc("GET /habits/{id}", controllers.GetHabitByID)
	mux.HandleFunc("PUT /habits/{id}", controllers.UpdateHabit)
	mux.HandleFunc("DELETE /habits/{id}", controllers.DeleteHabit)

	mux.HandleFunc("POST /habits/{id}/sessions", controllers.CreateSession)
	mux.HandleFunc("GET /habits/{id}/sessions", controllers.GetSessionsForHabit)
	mux.HandleFunc("PUT /sessions/{id}", controllers.UpdateSession)
	mux.HandleFunc("DELETE /sessions/{id}", controllers.DeleteSession)

}
