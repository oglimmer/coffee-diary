// Package service — Apple Sign in with Apple server-to-server token helpers.
// Used for obtaining refresh tokens at login and revoking them at account deletion,
// as required by App Store Review Guideline 5.1.1(v).
package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// AppleTokenService exchanges Apple authorization codes for refresh tokens
// and revokes them on account deletion. It is a no-op if server-to-server
// credentials (team ID, key ID, private key) are not configured.
type AppleTokenService struct {
	clientID   string
	teamID     string
	keyID      string
	privateKey *ecdsa.PrivateKey
	httpClient *http.Client
}

// NewAppleTokenService constructs the service. If any credential is missing,
// it returns a service with privateKey=nil that logs warnings and skips calls.
func NewAppleTokenService(clientID, teamID, keyID, privateKeyPEM string) *AppleTokenService {
	s := &AppleTokenService{
		clientID:   clientID,
		teamID:     teamID,
		keyID:      keyID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	if teamID == "" || keyID == "" || privateKeyPEM == "" {
		slog.Warn("Apple server-to-server credentials not fully configured — token exchange and revocation disabled")
		return s
	}
	key, err := parseECPrivateKey(privateKeyPEM)
	if err != nil {
		slog.Error("failed to parse Apple private key", "error", err)
		return s
	}
	s.privateKey = key
	return s
}

// Configured reports whether server-to-server calls can be made.
func (s *AppleTokenService) Configured() bool { return s.privateKey != nil }

// ExchangeCode trades an authorization code for tokens and returns the refresh token.
// Returns an empty string (and no error) if the service isn't configured.
func (s *AppleTokenService) ExchangeCode(ctx context.Context, code string) (string, error) {
	if !s.Configured() {
		return "", nil
	}
	clientSecret, err := s.clientSecret()
	if err != nil {
		return "", fmt.Errorf("apple client_secret: %w", err)
	}

	form := url.Values{}
	form.Set("client_id", s.clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://appleid.apple.com/auth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("apple token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}
	return tokenResp.RefreshToken, nil
}

// Revoke invalidates a refresh token at Apple. No-op if not configured or token empty.
func (s *AppleTokenService) Revoke(ctx context.Context, refreshToken string) error {
	if !s.Configured() || refreshToken == "" {
		return nil
	}
	clientSecret, err := s.clientSecret()
	if err != nil {
		return fmt.Errorf("apple client_secret: %w", err)
	}

	form := url.Values{}
	form.Set("client_id", s.clientID)
	form.Set("client_secret", clientSecret)
	form.Set("token", refreshToken)
	form.Set("token_type_hint", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://appleid.apple.com/auth/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("apple revoke endpoint returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// clientSecret builds the signed JWT used to authenticate to Apple's token endpoint.
// Valid for 10 minutes (Apple allows up to 6 months, but short-lived is safer).
func (s *AppleTokenService) clientSecret() (string, error) {
	now := time.Now()
	claims := jwt.Claims{
		Issuer:   s.teamID,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(10 * time.Minute)),
		Audience: jwt.Audience{"https://appleid.apple.com"},
		Subject:  s.clientID,
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: s.privateKey},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", s.keyID),
	)
	if err != nil {
		return "", err
	}
	return jwt.Signed(signer).Claims(claims).Serialize()
}

func parseECPrivateKey(pemStr string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid PEM data")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fall back to SEC1 format in case the key is not PKCS8-wrapped
		ecKey, ecErr := x509.ParseECPrivateKey(block.Bytes)
		if ecErr != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		return ecKey, nil
	}
	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not ECDSA")
	}
	return ecKey, nil
}
