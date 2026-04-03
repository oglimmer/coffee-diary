// Migrated from: CoffeeRepository.java
package repository

import (
	"context"
	"database/sql"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
)

type CoffeeRepository struct {
	db *sql.DB
}

func NewCoffeeRepository(db *sql.DB) *CoffeeRepository {
	return &CoffeeRepository{db: db}
}

func (r *CoffeeRepository) FindAllByUserID(ctx context.Context, userID int64) ([]domain.Coffee, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, user_id FROM coffees WHERE user_id = ?", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coffees []domain.Coffee
	for rows.Next() {
		var c domain.Coffee
		if err := rows.Scan(&c.ID, &c.Name, &c.UserID); err != nil {
			return nil, err
		}
		coffees = append(coffees, c)
	}
	return coffees, rows.Err()
}

func (r *CoffeeRepository) FindByID(ctx context.Context, id int64) (*domain.Coffee, error) {
	c := &domain.Coffee{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, user_id FROM coffees WHERE id = ?", id,
	).Scan(&c.ID, &c.Name, &c.UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CoffeeRepository) Create(ctx context.Context, userID int64, name string) (*domain.Coffee, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO coffees (name, user_id) VALUES (?, ?)", name, userID,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &domain.Coffee{ID: id, Name: name, UserID: userID}, nil
}

func (r *CoffeeRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM coffees WHERE id = ?", id)
	return err
}
