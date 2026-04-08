package handler

import (
	"encoding/json"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type FlagHandler struct {
	flagService *service.FlagService
}

func NewFlagHandler(flagService *service.FlagService) *FlagHandler {
	return &FlagHandler{flagService: flagService}
}

func (h *FlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID   string `json:"project_id"`
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

func (h *FlagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var req struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

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
