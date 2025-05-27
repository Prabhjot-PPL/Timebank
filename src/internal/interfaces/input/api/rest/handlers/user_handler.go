package userhandler

import (
	"context"
	"encoding/json"
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
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to register user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	err = u.userService.RegisterUser(ctx, user)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to register user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "success",
		Message: "User Registered Successfully ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

// LOGIN
func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var requestUser dto.UserDetails
	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to login user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	loginResponse, err := u.userService.LoginUser(ctx, requestUser)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to login user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
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

	response := pkg.StandardResponse{
		Status:  "success",
		Message: "User Logged-in Successfully ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)

}

func (u *UserHandler) FindHelperHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	skill := r.URL.Query().Get("skill")
	if skill == "" {
		http.Error(w, "Missing skill query param", http.StatusBadRequest)
		return
	}

	users, err := u.userService.FindUsersBySkill(ctx, skill)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to find helper) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "success",
		Data:    users,
		Message: "Query Successful ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

// // --------------------------------SESSION----------------------------------------

func (u *UserHandler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var session dto.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to create session) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	err := u.userService.CreateSession(ctx, session)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to create session) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "success",
		Message: "Session created Successful ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) StartSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract session ID from URL
	sessionIDStr := chi.URLParam(r, "id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to start session) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	err = u.userService.StartSession(ctx, sessionID)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to create session user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "success",
		Message: "Session started Successful ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) CompleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract session ID from URL
	sessionIDStr := chi.URLParam(r, "id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to start session) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	// Decode request body
	var req struct {
		Feedback string `json:"feedback"`
		Status   string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to start session) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := u.userService.CompleteSession(ctx, sessionID, req.Feedback, req.Status); err != nil {
		response := pkg.StandardResponse{
			Status:  "failure",
			Message: "(Failed to create session user) " + err.Error(),
		}
		pkg.WriteResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "success",
		Message: "Session completed Successful ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}
