package userhandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"timebank/src/internal/core/coreinterfaces"
	"timebank/src/internal/core/dto"
	"timebank/src/pkg"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService coreinterfaces.Service
}

func NewUserHandler(userService coreinterfaces.Service) coreinterfaces.UserAPIHandler {
	return &UserHandler{userService: userService}
}

// --------------------------------AUTH----------------------------------------

// REGISTER USER
func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var user dto.UserDetails

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error white decoding request : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = u.userService.RegisterUser(ctx, user)
	resp := "User Registered Successfully"

	if err != nil {
		log.Println("Error in registeration response : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	errString := "Failed to Register user."
	successString := "User Registered successfully"
	pkg.WriteResponse(w, resp, errString, successString)
}

// LOGIN
func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var requestUser dto.UserDetails
	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	loginResponse, err := u.userService.LoginUser(ctx, requestUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    loginResponse.TokenString,
		Expires:  loginResponse.TokenExpire,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &atCookie)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", loginResponse.FoundUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "successful login"})

}

// --------------------------------SESSION----------------------------------------

func (u *UserHandler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var session dto.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := u.userService.CreateSession(ctx, session)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Session created successfully"))
}

// func (u *UserHandler) CompleteSessionHandler(w http.ResponseWriter, r *http.Request) {
// }

func (u *UserHandler) CompleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract session ID from URL
	sessionIDStr := chi.URLParam(r, "id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req struct {
		Feedback string `json:"feedback"`
		Status   string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := u.userService.CompleteSession(ctx, sessionID, req.Feedback, req.Status); err != nil {
		http.Error(w, "Failed to complete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session completed successfully"))
}
