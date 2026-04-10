package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest  true  "Register data"
// @Success      201      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/register [post]
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
	tokenPair, err := a.authService.Register(r.Context(), req.Email, req.Password, orgID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]any{"access_token": tokenPair.AccessToken, "refresh_token": tokenPair.RefreshToken})
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and get token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "Login data"
// @Success      200      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/login [post]
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
	respondJSON(w, http.StatusOK, map[string]any{"access_token": tokenPair.AccessToken, "refresh_token": tokenPair.RefreshToken})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshRequest  true  "Refresh token"
// @Success      200      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/refresh [post]
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

// Logout godoc
// @Summary      Logout
// @Description  Invalidate refresh token
// @Tags         auth
// @Accept       json
// @Param        request  body  logoutRequest  true  "Refresh token"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/logout [post]
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
