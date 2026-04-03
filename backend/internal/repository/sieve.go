// Migrated from: SieveRepository.java
package repository

import (
	"context"
	"database/sql"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
)

type SieveRepository struct {
	db *sql.DB
}

func NewSieveRepository(db *sql.DB) *SieveRepository {
	return &SieveRepository{db: db}
}

func (r *SieveRepository) FindAllByUserID(ctx context.Context, userID int64) ([]domain.Sieve, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, user_id FROM sieves WHERE user_id = ?", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sieves []domain.Sieve
	for rows.Next() {
		var s domain.Sieve
		if err := rows.Scan(&s.ID, &s.Name, &s.UserID); err != nil {
			return nil, err
		}
		sieves = append(sieves, s)
	}
	return sieves, rows.Err()
}

func (r *SieveRepository) FindByID(ctx context.Context, id int64) (*domain.Sieve, error) {
	s := &domain.Sieve{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, user_id FROM sieves WHERE id = ?", id,
	).Scan(&s.ID, &s.Name, &s.UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *SieveRepository) Create(ctx context.Context, userID int64, name string) (*domain.Sieve, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO sieves (name, user_id) VALUES (?, ?)", name, userID,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &domain.Sieve{ID: id, Name: name, UserID: userID}, nil
}

func (r *SieveRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sieves WHERE id = ?", id)
	return err
}
