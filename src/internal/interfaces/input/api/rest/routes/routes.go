package routes

import (
	"net/http"
	"timebank/src/internal/core/coreinterfaces"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler coreinterfaces.UserAPIHandler) http.Handler {
	r := chi.NewRouter()

	// REGISTER
	r.Post("/auth/register", userHandler.RegisterHandler)

	// LOGIN
	r.Post("/auth/login", userHandler.LoginHandler)

	// VIEW PROFILE
	// r.Get("/users/me", userHandler.ViewProfileHandler)

	// UPDATE PROFILE
	// r.Put("/users/me", userHandler.UpdateProfileHandler)

	// DELETE PROFILE
	// r.Delete("/users/me", userHandler.DeleteProfileHandler)

	// FIND HELPER
	// r.Get("/users/search?skill=guitar", userHandler.FindHelperHandler)

	// CREATE SESSION
	r.Post("/sessions", userHandler.CreateSessionHandler)

	// START SESSION
	// r.Post("/sessions/55/start", userHandler.StartSessionHandler)

	// END SESSION
	r.Post("/sessions/{id}/complete", userHandler.CompleteSessionHandler)

	return r
}
