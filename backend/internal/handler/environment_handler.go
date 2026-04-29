package handler

import (
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/port"
	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	envService port.EnvironmentServiceInterface
}

func NewEnvironmentHandler(envService port.EnvironmentServiceInterface) *EnvironmentHandler {
	return &EnvironmentHandler{envService: envService}
}

func (h *EnvironmentHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(r.URL.Query().Get("project_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project_id")
		return
	}
	envs, err := h.envService.ListEnvironments(r.Context(), projectID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondJSON(w, http.StatusOK, envs)
}
