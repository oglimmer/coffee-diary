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
	appleTokens        *service.AppleTokenService
	oauth2Config       *oauth2.Config
	verifier           *oidc.IDTokenVerifier
	appleVerifier      *oidc.IDTokenVerifier
	store              sessions.Store
	frontendURL        string
	endSessionEndpoint string
}

func NewAuthHandler(
	authService *service.AuthService,
	appleTokens *service.AppleTokenService,
	oauth2Config *oauth2.Config,
	verifier *oidc.IDTokenVerifier,
	appleVerifier *oidc.IDTokenVerifier,
	store sessions.Store,
	frontendURL string,
	endSessionEndpoint string,
) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		appleTokens:        appleTokens,
		oauth2Config:       oauth2Config,
		verifier:           verifier,
		appleVerifier:      appleVerifier,
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

// AppleCallback handles Sign in with Apple token verification.
func (h *AuthHandler) AppleCallback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdentityToken     string `json:"identityToken"`
		AuthorizationCode string `json:"authorizationCode"`
		FullName          string `json:"fullName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.WriteError(w, apperr.BadRequest("Invalid request body"))
		return
	}

	if req.IdentityToken == "" {
		apperr.WriteError(w, apperr.BadRequest("identityToken is required"))
		return
	}

	// Verify Apple identity token using OIDC discovery
	idToken, err := h.appleVerifier.Verify(r.Context(), req.IdentityToken)
	if err != nil {
		slog.Error("Apple token verification failed", "error", err)
		apperr.WriteError(w, apperr.Unauthorized("Authentication failed"))
		return
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("failed to parse Apple claims", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	username := req.FullName
	if username == "" {
		username = claims.Email
	}
	if username == "" {
		username = claims.Sub
	}

	// Prefix sub to distinguish Apple users from Keycloak users
	appleSub := "apple:" + claims.Sub

	user, err := h.authService.FindOrCreateByOIDC(r.Context(), appleSub, username)
	if err != nil {
		slog.Error("failed to find/create user", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	// Exchange the authorization code for a refresh token so we can revoke it
	// when the user deletes their account (App Store Guideline 5.1.1(v)).
	if req.AuthorizationCode != "" && h.appleTokens.Configured() {
		refreshToken, err := h.appleTokens.ExchangeCode(r.Context(), req.AuthorizationCode)
		if err != nil {
			slog.Warn("Apple code exchange failed — account deletion revocation will be skipped",
				"userID", user.ID, "error", err)
		} else if refreshToken != "" {
			if err := h.authService.StoreAppleRefreshToken(r.Context(), user.ID, refreshToken); err != nil {
				slog.Warn("failed to persist Apple refresh token", "userID", user.ID, "error", err)
			}
		}
	}

	sess, _ := h.store.Get(r, sessionName)
	sess.Values["userID"] = user.ID
	if err := sess.Save(r, w); err != nil {
		slog.Error("failed to save session", "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	slog.Info("user logged in via Apple", "username", user.Username, "id", user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	})
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

// DeleteAccount permanently removes the authenticated user's account and all their data.
// Required by App Store Review Guideline 5.1.1(v). For Sign in with Apple users this
// also revokes the refresh token at Apple.
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, sessionName)
	if err != nil {
		apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
		return
	}
	uid, ok := sess.Values["userID"].(int64)
	if !ok {
		apperr.WriteError(w, apperr.Unauthorized("Authentication required"))
		return
	}

	if err := h.authService.DeleteAccount(r.Context(), uid); err != nil {
		slog.Error("failed to delete account", "userID", uid, "error", err)
		apperr.WriteError(w, apperr.InternalError())
		return
	}

	// Invalidate the session cookie.
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		slog.Warn("failed to clear session after account deletion", "error", err)
	}

	slog.Info("user account deleted", "userID", uid)
	w.WriteHeader(http.StatusNoContent)
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
