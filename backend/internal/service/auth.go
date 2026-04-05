// Migrated from: AuthService.java + AppUserDetailsService.java
package service

import (
	"context"
	"log/slog"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	appleTokens *AppleTokenService
}

func NewAuthService(userRepo *repository.UserRepository, appleTokens *AppleTokenService) *AuthService {
	return &AuthService{userRepo: userRepo, appleTokens: appleTokens}
}

// StoreAppleRefreshToken persists a refresh token obtained from Apple's token endpoint.
func (s *AuthService) StoreAppleRefreshToken(ctx context.Context, userID int64, token string) error {
	return s.userRepo.SetAppleRefreshToken(ctx, userID, token)
}

// DeleteAccount removes all user data and, for Apple sign-in users, revokes
// their refresh token at Apple as required by App Store Review Guideline 5.1.1(v).
func (s *AuthService) DeleteAccount(ctx context.Context, userID int64) error {
	// Revoke Apple tokens first — if this fails we still proceed with local deletion
	// (better to leave a stale Apple token than to leave the user's data behind).
	refreshToken, err := s.userRepo.GetAppleRefreshToken(ctx, userID)
	if err != nil {
		slog.Warn("failed to read Apple refresh token for revocation", "userID", userID, "error", err)
	} else if refreshToken != "" {
		if err := s.appleTokens.Revoke(ctx, refreshToken); err != nil {
			slog.Error("failed to revoke Apple refresh token", "userID", userID, "error", err)
		}
	}

	return s.userRepo.DeleteUserCascade(ctx, userID)
}

// FindByID returns a user by their database ID.
func (s *AuthService) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// FindOrCreateByOIDC looks up a user by OIDC sub claim. If not found, tries to
// match by username (for migrating existing users) or creates a new user.
func (s *AuthService) FindOrCreateByOIDC(ctx context.Context, sub, username string) (*domain.User, error) {
	// 1. Try by OIDC sub
	user, err := s.userRepo.FindByOIDCSub(ctx, sub)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 2. Try by username (existing user migrating to OIDC)
	user, err = s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := s.userRepo.UpdateOIDCSub(ctx, user.ID, sub); err != nil {
			return nil, err
		}
		user.OIDCSub = sub
		return user, nil
	}

	// 3. Create new user
	return s.userRepo.Create(ctx, username, sub)
}
