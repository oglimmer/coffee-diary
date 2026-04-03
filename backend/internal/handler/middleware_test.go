package handler

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func init() {
	gob.Register(int64(0))
}

func TestSecurityHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(inner)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	assert.Equal(t, "default-src 'self'", rec.Header().Get("Content-Security-Policy"))
}

func TestSessionAuth_NoSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("test-secret-32-chars-long-enough"))

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SessionAuth(store)(inner)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/coffees", nil)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSessionAuth_WithSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("test-secret-32-chars-long-enough"))

	// Create a session
	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	sess, _ := store.Get(req1, sessionName)
	sess.Values["userID"] = int64(42)
	sess.Save(req1, rec1)

	// Extract cookie
	cookies := rec1.Result().Cookies()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := UserIDFromContext(r.Context())
		assert.Equal(t, int64(42), uid)
		w.WriteHeader(http.StatusOK)
	})

	handler := SessionAuth(store)(inner)
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/coffees", nil)
	for _, c := range cookies {
		req2.AddCookie(c)
	}
	handler.ServeHTTP(rec2, req2)

	assert.Equal(t, http.StatusOK, rec2.Code)
}
