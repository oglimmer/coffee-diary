// Migrated from: UserRepository.java
package repository

import (
	"context"
	"database/sql"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	u := &domain.User{}
	var oidcSub sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, oidc_sub, created_at FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Username, &oidcSub, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.OIDCSub = oidcSub.String
	return u, err
}

// GetAppleRefreshToken returns the stored Apple refresh token for a user, or "" if none.
func (r *UserRepository) GetAppleRefreshToken(ctx context.Context, id int64) (string, error) {
	var token sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT apple_refresh_token FROM users WHERE id = ?", id,
	).Scan(&token)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return token.String, nil
}

// SetAppleRefreshToken stores (or clears) the Apple refresh token for a user.
func (r *UserRepository) SetAppleRefreshToken(ctx context.Context, id int64, token string) error {
	var val sql.NullString
	if token != "" {
		val = sql.NullString{String: token, Valid: true}
	}
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET apple_refresh_token = ? WHERE id = ?", val, id,
	)
	return err
}

// DeleteUserCascade removes the user and all owned data in a single transaction.
func (r *UserRepository) DeleteUserCascade(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		"DELETE FROM diary_entries WHERE user_id = ?",
		"DELETE FROM coffees WHERE user_id = ?",
		"DELETE FROM sieves WHERE user_id = ?",
		"DELETE FROM users WHERE id = ?",
	}
	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *UserRepository) FindByOIDCSub(ctx context.Context, sub string) (*domain.User, error) {
	u := &domain.User{}
	var oidcSub sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, oidc_sub, created_at FROM users WHERE oidc_sub = ?", sub,
	).Scan(&u.ID, &u.Username, &oidcSub, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.OIDCSub = oidcSub.String
	return u, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	u := &domain.User{}
	var oidcSub sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, oidc_sub, created_at FROM users WHERE username = ?", username,
	).Scan(&u.ID, &u.Username, &oidcSub, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	u.OIDCSub = oidcSub.String
	return u, err
}

func (r *UserRepository) Create(ctx context.Context, username, oidcSub string) (*domain.User, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (username, oidc_sub) VALUES (?, ?)", username, oidcSub,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: id, Username: username, OIDCSub: oidcSub}, nil
}

func (r *UserRepository) UpdateOIDCSub(ctx context.Context, id int64, oidcSub string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET oidc_sub = ? WHERE id = ?", oidcSub, id,
	)
	return err
}
