package main

import (
	"log"
	"net/http"

	"github.com/mehradSaeedi/habit-tracker/pkg/config"
	"github.com/mehradSaeedi/habit-tracker/pkg/models"
	"github.com/mehradSaeedi/habit-tracker/pkg/routes"
)

func main() {
	config.Connect()

	config.DB.AutoMigrate(
		&models.Habit{},
		&models.Session{},
	)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
