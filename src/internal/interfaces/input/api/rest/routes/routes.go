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

	// FIND HELPER
	r.Get("/users/search", userHandler.FindHelperHandler)

	// CREATE SESSION
	r.Post("/sessions", userHandler.CreateSessionHandler)

	// // START SESSION
	r.Put("/session/{id}/start", userHandler.StartSessionHandler)

	// // END SESSION
	r.Post("/sessions/{id}/complete", userHandler.CompleteSessionHandler)

	return r
}
