// Migrated from: Spring Boot Actuator endpoints
package handler

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/oglimmer/coffee-diary-backend/internal/config"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
)

type ActuatorHandler struct {
	cfg *config.Config
}

func NewActuatorHandler(cfg *config.Config) *ActuatorHandler {
	return &ActuatorHandler{cfg: cfg}
}

func (h *ActuatorHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func (h *ActuatorHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"app": map[string]string{
			"name":    h.cfg.AppName,
			"version": h.cfg.AppVersion,
		},
		"build": map[string]string{
			"time": h.cfg.BuildTime,
		},
		"git": map[string]string{
			"commit": h.cfg.GitCommit,
		},
	})
}

func (h *ActuatorHandler) Prometheus() http.Handler {
	return h.basicAuth(promhttp.Handler())
}

func (h *ActuatorHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	// Simple metrics endpoint — Prometheus handler covers detailed metrics
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "available"})
}

func (h *ActuatorHandler) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(h.cfg.ActuatorUsername)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(h.cfg.ActuatorPassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Actuator"`)
			apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
