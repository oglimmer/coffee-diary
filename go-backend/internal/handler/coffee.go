// Migrated from: CoffeeController.java
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

type CoffeeHandler struct {
	coffeeService *service.CoffeeService
}

func NewCoffeeHandler(coffeeService *service.CoffeeService) *CoffeeHandler {
	return &CoffeeHandler{coffeeService: coffeeService}
}

func (h *CoffeeHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	coffees, err := h.coffeeService.FindAllByUser(r.Context(), userID)
	if err != nil {
		slog.Error("failed to find coffees", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffees)
}

func (h *CoffeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	var req domain.CoffeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid request body"))
		return
	}

	resp, err := h.coffeeService.Create(r.Context(), userID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *CoffeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid ID"))
		return
	}

	if err := h.coffeeService.Delete(r.Context(), userID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleServiceError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperr.AppError); ok {
		apperr.WriteError(w, appErr)
		return
	}
	slog.Error("unexpected error", "error", err)
	apperr.WriteError(w, apperr.InternalError())
}
