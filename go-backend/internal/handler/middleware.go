// Migrated from: SecurityConfig.java
package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"

	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
)

type contextKey string

const userIDKey contextKey = "userID"
const sessionName = "session"

// SessionAuth middleware checks for a valid session with a user ID.
func SessionAuth(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := store.Get(r, sessionName)
			if err != nil {
				apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
				return
			}

			uid, ok := sess.Values["userID"]
			if !ok {
				apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
				return
			}

			userID, ok := uid.(int64)
			if !ok {
				apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the authenticated user ID from context.
func UserIDFromContext(ctx context.Context) int64 {
	return ctx.Value(userIDKey).(int64)
}

// SecurityHeaders adds standard security headers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}
