package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SDKKeyHandler struct {
	repo domain.SDKKeyRepository
	db   *sqlx.DB
}

func NewSDKKeyHandler(repo domain.SDKKeyRepository) *SDKKeyHandler {
	return &SDKKeyHandler{repo: repo}
}

func (h *SDKKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID     string `json:"project_id"`
		EnvironmentID string `json:"environment_id"`
		Name          string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project_id")
		return
	}

	envID, err := uuid.Parse(req.EnvironmentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid environment_id")
		return
	}

	actorID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	rawKey := "sdk-" + hex.EncodeToString(rawBytes)

	key := &domain.SDKKey{
		ProjectID:     projectID,
		EnvironmentID: envID,
		Name:          req.Name,
	}
	if actorID != uuid.Nil {
		key.CreatedBy = actorID
	}

	if err := h.repo.Create(r.Context(), key, rawKey); err != nil {
		slog.Error("sdk key create failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"id":  key.ID,
		"key": rawKey,
	})
}

func (h *SDKKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(r.URL.Query().Get("project_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project_id")
		return
	}

	keys, err := h.repo.ListByProject(r.Context(), projectID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if keys == nil {
		keys = []*domain.SDKKey{}
	}
	respondJSON(w, http.StatusOK, keys)
}

func (h *SDKKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		domainErrorToHTTP(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
