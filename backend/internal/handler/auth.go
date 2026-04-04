// Migrated from: AuthController.java (rewritten for OIDC)
package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/service"
)

type AuthHandler struct {
	authService        *service.AuthService
	oauth2Config       *oauth2.Config
	verifier           *oidc.IDTokenVerifier
	store              sessions.Store
	frontendURL        string
	endSessionEndpoint string
}

func NewAuthHandler(
	authService *service.AuthService,
	oauth2Config *oauth2.Config,
	verifier *oidc.IDTokenVerifier,
	store sessions.Store,
	frontendURL string,
	endSessionEndpoint string,
) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		oauth2Config:       oauth2Config,
		verifier:           verifier,
		store:              store,
		frontendURL:        frontendURL,
		endSessionEndpoint: endSessionEndpoint,
	}
}

// Login redirects to the OIDC provider.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		slog.Error("failed to generate state", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	sess, _ := h.store.Get(r, sessionName)
	sess.Values["oauth_state"] = state

	// Allow mobile clients to specify a custom redirect after login (e.g. coffeeDiary://auth/callback)
	if redirectAfter := r.URL.Query().Get("redirect_after"); redirectAfter != "" {
		sess.Values["redirect_after"] = redirectAfter
	}

	if err := sess.Save(r, w); err != nil {
		slog.Error("failed to save session", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	http.Redirect(w, r, h.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

// Callback handles the OIDC provider redirect.
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.store.Get(r, sessionName)

	// Verify state
	expectedState, _ := sess.Values["oauth_state"].(string)
	if r.URL.Query().Get("state") != expectedState || expectedState == "" {
		slog.Warn("OIDC state mismatch")
		apperr.WriteError(w, apperr.BadRequest("Invalid state parameter"))
		return
	}
	delete(sess.Values, "oauth_state")

	// Exchange code for token
	oauth2Token, err := h.oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		slog.Error("OIDC token exchange failed", "error", err)
		apperr.WriteError(w, apperr.Unauthorized("Authentication failed"))
		return
	}

	// Extract and verify ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		slog.Error("no id_token in OIDC response")
		apperr.WriteError(w, apperr.Unauthorized("Authentication failed"))
		return
	}

	idToken, err := h.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		slog.Error("OIDC token verification failed", "error", err)
		apperr.WriteError(w, apperr.Unauthorized("Authentication failed"))
		return
	}

	var claims struct {
		Sub               string `json:"sub"`
		PreferredUsername  string `json:"preferred_username"`
		Email             string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("failed to parse OIDC claims", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	username := claims.PreferredUsername
	if username == "" {
		username = claims.Email
	}

	// Find or create user
	user, err := h.authService.FindOrCreateByOIDC(r.Context(), claims.Sub, username)
	if err != nil {
		slog.Error("failed to find/create user", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	// Extract mobile redirect before overwriting session values
	redirectAfter, _ := sess.Values["redirect_after"].(string)
	delete(sess.Values, "redirect_after")

	sess.Values["userID"] = user.ID
	sess.Values["id_token"] = rawIDToken
	if err := sess.Save(r, w); err != nil {
		slog.Error("failed to save session", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	slog.Info("user logged in via OIDC", "username", user.Username, "id", user.ID)

	// For mobile clients: redirect to custom scheme with session cookie value
	if redirectAfter != "" && strings.HasPrefix(redirectAfter, "coffeeDiary://") {
		var sessionCookie string
		for _, c := range w.Header()["Set-Cookie"] {
			if strings.HasPrefix(c, sessionName+"=") {
				parts := strings.SplitN(c, ";", 2)
				sessionCookie = strings.TrimPrefix(parts[0], sessionName+"=")
				break
			}
		}
		redirectURL := redirectAfter + "?session_cookie=" + url.QueryEscape(sessionCookie)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	http.Redirect(w, r, h.frontendURL, http.StatusFound)
}

// Logout clears the local session and redirects to Keycloak's end_session_endpoint.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.store.Get(r, sessionName)
	idToken, _ := sess.Values["id_token"].(string)

	sess.Options.MaxAge = -1
	sess.Save(r, w)

	if h.endSessionEndpoint != "" {
		logoutURL, err := url.Parse(h.endSessionEndpoint)
		if err == nil {
			q := logoutURL.Query()
			if idToken != "" {
				q.Set("id_token_hint", idToken)
			}
			q.Set("client_id", h.oauth2Config.ClientID)
			q.Set("post_logout_redirect_uri", h.frontendURL)
			logoutURL.RawQuery = q.Encode()
			http.Redirect(w, r, logoutURL.String(), http.StatusFound)
			return
		}
		slog.Error("failed to parse end_session_endpoint", "error", err)
	}

	http.Redirect(w, r, h.frontendURL, http.StatusFound)
}

// Me returns the current authenticated user.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, sessionName)
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

	// Look up user to get username
	user, err := h.authService.FindByID(r.Context(), userID)
	if err != nil || user == nil {
		apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
