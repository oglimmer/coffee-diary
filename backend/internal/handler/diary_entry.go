// Migrated from: DiaryEntryController.java
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
	"github.com/oglimmer/coffee-diary-backend/internal/service"
)

type DiaryEntryHandler struct {
	entryService *service.DiaryEntryService
}

func NewDiaryEntryHandler(entryService *service.DiaryEntryService) *DiaryEntryHandler {
	return &DiaryEntryHandler{entryService: entryService}
}

func (h *DiaryEntryHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	q := r.URL.Query()

	filter := repository.DiaryEntryFilter{UserID: userID}

	if v := q.Get("coffeeId"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			filter.CoffeeID = &id
		}
	}
	if v := q.Get("sieveId"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			filter.SieveID = &id
		}
	}
	if v := q.Get("dateFrom"); v != "" {
		t, err := time.Parse("2006-01-02T15:04:05", v)
		if err == nil {
			filter.DateFrom = &t
		}
	}
	if v := q.Get("dateTo"); v != "" {
		t, err := time.Parse("2006-01-02T15:04:05", v)
		if err == nil {
			filter.DateTo = &t
		}
	}
	if v := q.Get("ratingMin"); v != "" {
		rating, err := strconv.Atoi(v)
		if err == nil {
			filter.RatingMin = &rating
		}
	}

	page := 0
	if v := q.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		}
	}
	size := 20
	if v := q.Get("size"); v != "" {
		if s, err := strconv.Atoi(v); err == nil && s > 0 {
			size = s
		}
	}

	sortField := "dateTime"
	sortDir := "asc"
	if v := q.Get("sort"); v != "" {
		parts := strings.SplitN(v, ",", 2)
		sortField = parts[0]
		if len(parts) > 1 {
			sortDir = parts[1]
		}
	}

	resp, err := h.entryService.FindAll(r.Context(), filter, page, size, sortField, sortDir)
	if err != nil {
		slog.Error("failed to find diary entries", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DiaryEntryHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid ID"))
		return
	}

	resp, err := h.entryService.FindByID(r.Context(), userID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DiaryEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	var req domain.DiaryEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid request body"))
		return
	}

	if req.DateTime.Time.IsZero() {
		apperr.WriteValidationError(w, map[string]string{"dateTime": "Date/time is required"})
		return
	}
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		apperr.WriteValidationError(w, map[string]string{"rating": "Rating must be between 1 and 5"})
		return
	}

	resp, err := h.entryService.Create(r.Context(), userID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *DiaryEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid ID"))
		return
	}

	var req domain.DiaryEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid request body"))
		return
	}

	if req.DateTime.Time.IsZero() {
		apperr.WriteValidationError(w, map[string]string{"dateTime": "Date/time is required"})
		return
	}
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		apperr.WriteValidationError(w, map[string]string{"rating": "Rating must be between 1 and 5"})
		return
	}

	resp, err := h.entryService.Update(r.Context(), userID, id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DiaryEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid ID"))
		return
	}

	if err := h.entryService.Delete(r.Context(), userID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
