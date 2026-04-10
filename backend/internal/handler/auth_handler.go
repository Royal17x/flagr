package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	orgID, err := uuid.Parse(req.OrgID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := a.authService.Register(r.Context(), req.Email, req.Password, orgID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]any{"user_id": id})
}
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tokenPair, err := a.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"acess_token": tokenPair.AccessToken, "refresh_token": tokenPair.RefreshToken})
}
func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tokenPair, err := a.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"access_token": tokenPair.AccessToken, "refresh_token": tokenPair.RefreshToken})
}
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := a.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
