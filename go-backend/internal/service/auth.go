// Migrated from: AuthService.java + AppUserDetailsService.java
package service

import (
	"context"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
