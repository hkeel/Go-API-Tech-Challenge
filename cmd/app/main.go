package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hkeel/Go-API-Tech-Challenge/internal/config"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/database"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/handlers"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/repositories"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Starting application...")
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database initialized successfully")

	// Initialize services
	courseRepo := &repositories.CourseRepository{DB: db}
	courseService := &services.CourseService{Repo: courseRepo}

	// Initialize handlers
	courseHandler := handlers.NewCourseHandler(courseService)
	//personHandler := handlers.NewPersonHandler(personService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(api chi.Router) {
		// Health check endpoint
		api.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
		})

		/// Course routes
		api.Route("/course", func(r chi.Router) {
			r.Get("/", courseHandler.GetAllCourses)
		})

		// Person routes
		// api.Route("/person", func(r chi.Router) {
		// 	r.Get("/", personHandler.GetAllPeople)
		// 	// Add more person routes here
		// })
	})

	http.ListenAndServe(":8000", r)
}
