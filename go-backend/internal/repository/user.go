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
