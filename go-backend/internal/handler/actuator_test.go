package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oglimmer/coffee-diary-backend/internal/config"
)

func TestActuator_Health(t *testing.T) {
	h := NewActuatorHandler(&config.Config{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var m map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m))
	assert.Equal(t, "UP", m["status"])
}

func TestActuator_Info(t *testing.T) {
	h := NewActuatorHandler(&config.Config{AppName: "test-app", AppVersion: "1.0.0"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/actuator/info", nil)
	h.Info(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m))
	app := m["app"].(map[string]interface{})
	assert.Equal(t, "test-app", app["name"])
	assert.Equal(t, "1.0.0", app["version"])
}

func TestActuator_Prometheus_NoAuth(t *testing.T) {
	h := NewActuatorHandler(&config.Config{ActuatorUsername: "admin", ActuatorPassword: "secret"})
	handler := h.Prometheus()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/actuator/prometheus", nil)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestActuator_Prometheus_WithAuth(t *testing.T) {
	h := NewActuatorHandler(&config.Config{ActuatorUsername: "admin", ActuatorPassword: "secret"})
	handler := h.Prometheus()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/actuator/prometheus", nil)
	req.SetBasicAuth("admin", "secret")
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
