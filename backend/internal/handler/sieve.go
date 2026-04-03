// Migrated from: SieveController.java
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/service"
)

type SieveHandler struct {
	sieveService *service.SieveService
}

func NewSieveHandler(sieveService *service.SieveService) *SieveHandler {
	return &SieveHandler{sieveService: sieveService}
}

func (h *SieveHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	sieves, err := h.sieveService.FindAllByUser(r.Context(), userID)
	if err != nil {
		slog.Error("failed to find sieves", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sieves)
}

func (h *SieveHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	var req domain.SieveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid request body"))
		return
	}

	resp, err := h.sieveService.Create(r.Context(), userID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *SieveHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid ID"))
		return
	}

	if err := h.sieveService.Delete(r.Context(), userID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
