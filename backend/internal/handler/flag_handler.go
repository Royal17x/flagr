package handler

import (
	"encoding/json"
	"github.com/Royal17x/flagr/backend/internal/validator"
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FlagHandler struct {
	flagService port.FlagServiceInterface
}

func NewFlagHandler(flagService port.FlagServiceInterface) *FlagHandler {
	return &FlagHandler{flagService: flagService}
}

// Create godoc
// @Summary      Create a feature flag
// @Description  Creates a new feature flag in a project
// @Tags         flags
// @Accept       json
// @Produce      json
// @Param        request  body      createFlagRequest  true  "Flag data"
// @Success      201      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security BearerAuth
// @Router       /flags [post]
func (h *FlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateFlagKey(req.Key); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	newFlag := domain.Flag{
		ProjectID:   projectUUID,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Type:        domain.FlagType(req.Type),
	}

	flagID, err := h.flagService.CreateFlag(r.Context(), &newFlag)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]any{"flag_id": flagID})
}

// GetByID godoc
// @Summary      Get flag by ID
// @Description  Returns a single feature flag
// @Tags         flags
// @Produce      json
// @Param        id   path      string  true  "Flag UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /flags/{id} [get]
func (h *FlagHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	flagID, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	gotFlag, err := h.flagService.GetFlag(r.Context(), flagID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"flag": *gotFlag})
}

// List godoc
// @Summary      Get a list of flags
// @Description  Show a current list of flags in project
// @Tags         flags
// @Produce      json
// @Param   project_id  query  string  true  "Project UUID"
// @Success      200      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security BearerAuth
// @Router       /flags [get]
func (h *FlagHandler) List(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("project_id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	projectID, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	flags, err := h.flagService.ListFlags(r.Context(), projectID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"flags": flags})
}

// Update godoc
// @Summary      Update flag
// @Description  Changes past flag fields on new fields
// @Tags         flags
// @Accept       json
// @Produce      json
// @Param  request  body  updateFlagRequest  true  "Update data"
// @Param        id   path      string  true  "Flag UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /flags/{id} [put]
func (h *FlagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var req updateFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	flagID, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	flag := domain.Flag{
		ID:          flagID,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Type:        domain.FlagType(req.Type),
	}
	err = h.flagService.UpdateFlag(r.Context(), &flag)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"flag": flag})
}

// Delete godoc
// @Summary      Delete flag
// @Description  Deletes flag with ID in param
// @Tags         flags
// @Param        id   path      string  true  "Flag UUID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /flags/{id} [delete]
func (h *FlagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	flagID, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	err = h.flagService.DeleteFlag(r.Context(), flagID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Evaluate godoc
// @Summary      Get "enabled" of flag
// @Description  Returns a state with current flag ID
// @Tags         flags
// @Produce      json
// @Param   key             query  string  true  "Flag key"
// @Param   project_id      query  string  true  "Project UUID"
// @Param   environment_id  query  string  true  "Environment UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security SDKKeyAuth
// @Router  /flags/evaluate [get]
func (h *FlagHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	environmentID := r.URL.Query().Get("environment_id")
	if environmentID == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	environmentUUID, err := uuid.Parse(environmentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	enabled, err := h.flagService.EvaluateFlag(r.Context(), key, projectUUID, environmentUUID)
	if err != nil {
		domainErrorToHTTP(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"enabled": enabled})
}
