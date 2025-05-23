package coreinterfaces

import "net/http"

type UserAPIHandler interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	CreateSessionHandler(w http.ResponseWriter, r *http.Request)
	CompleteSessionHandler(w http.ResponseWriter, r *http.Request)
}
