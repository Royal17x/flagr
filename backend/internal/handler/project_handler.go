package handler

import (
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/Royal17x/flagr/backend/internal/port"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectService port.ProjectServiceInterface
}

func NewProjectHandler(projectService port.ProjectServiceInterface) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := r.Context().Value(middleware.OrgIDKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	projects, err := h.projectService.ListProjects(r.Context(), orgID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondJSON(w, http.StatusOK, projects)
}
