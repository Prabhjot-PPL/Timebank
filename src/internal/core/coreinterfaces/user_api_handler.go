package coreinterfaces

import "net/http"

type UserAPIHandler interface {
	// AUTH
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	FindHelperHandler(w http.ResponseWriter, r *http.Request)
	// // SESSION
	CreateSessionHandler(w http.ResponseWriter, r *http.Request)
	StartSessionHandler(w http.ResponseWriter, r *http.Request)
	CompleteSessionHandler(w http.ResponseWriter, r *http.Request)
}
