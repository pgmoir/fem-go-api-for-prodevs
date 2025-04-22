package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/pgmoir/femGoProject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	// r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	// r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	return r
}
