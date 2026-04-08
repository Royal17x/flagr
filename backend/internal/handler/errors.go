package handler

import (
	"errors"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"net/http"
)

func domainErrorToHTTP(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		respondError(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		respondError(w, http.StatusConflict, "already exists")
	case errors.Is(err, domain.ErrInvalidInput):
		respondError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, domain.ErrUnauthorized):
		respondError(w, http.StatusUnauthorized, "unauthorized")
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
