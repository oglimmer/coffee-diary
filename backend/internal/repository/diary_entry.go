// Migrated from: DiaryEntryRepository.java + DiaryEntrySpecification.java
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
)

type DiaryEntryRepository struct {
	db *sql.DB
}

func NewDiaryEntryRepository(db *sql.DB) *DiaryEntryRepository {
	return &DiaryEntryRepository{db: db}
}

type DiaryEntryFilter struct {
	UserID    int64
	CoffeeID *int64
	SieveID  *int64
	DateFrom *time.Time
	DateTo   *time.Time
	RatingMin *int
}

// allowedSortColumns maps JSON field names to DB column names.
var allowedSortColumns = map[string]string{
	"dateTime":    "de.date_time",
	"temperature": "de.temperature",
	"rating":      "de.rating",
	"id":          "de.id",
}

func (r *DiaryEntryRepository) FindAll(ctx context.Context, filter DiaryEntryFilter, page, size int, sortField, sortDir string) ([]domain.DiaryEntry, int64, error) {
	where, args := buildWhere(filter)

	// Count
	var total int64
	countQuery := "SELECT COUNT(*) FROM diary_entries de " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Validate sort
	col, ok := allowedSortColumns[sortField]
	if !ok {
		col = "de.date_time"
	}
	dir := "ASC"
	if strings.EqualFold(sortDir, "desc") {
		dir = "DESC"
	}

	query := fmt.Sprintf(`SELECT de.id, de.user_id, de.date_time, de.sieve_id, s.name,
		de.temperature, de.coffee_id, c.name, de.grind_size, de.input_weight,
		de.output_weight, de.time_seconds, de.rating, de.notes
		FROM diary_entries de
		LEFT JOIN sieves s ON de.sieve_id = s.id
		LEFT JOIN coffees c ON de.coffee_id = c.id
		%s ORDER BY %s %s LIMIT ? OFFSET ?`, where, col, dir)

	args = append(args, size, page*size)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []domain.DiaryEntry
	for rows.Next() {
		var e domain.DiaryEntry
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.DateTime, &e.SieveID, &e.SieveName,
			&e.Temperature, &e.CoffeeID, &e.CoffeeName, &e.GrindSize, &e.InputWeight,
			&e.OutputWeight, &e.TimeSeconds, &e.Rating, &e.Notes,
		); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

func (r *DiaryEntryRepository) FindByID(ctx context.Context, id int64) (*domain.DiaryEntry, error) {
	e := &domain.DiaryEntry{}
	err := r.db.QueryRowContext(ctx, `SELECT de.id, de.user_id, de.date_time, de.sieve_id, s.name,
		de.temperature, de.coffee_id, c.name, de.grind_size, de.input_weight,
		de.output_weight, de.time_seconds, de.rating, de.notes
		FROM diary_entries de
		LEFT JOIN sieves s ON de.sieve_id = s.id
		LEFT JOIN coffees c ON de.coffee_id = c.id
		WHERE de.id = ?`, id,
	).Scan(
		&e.ID, &e.UserID, &e.DateTime, &e.SieveID, &e.SieveName,
		&e.Temperature, &e.CoffeeID, &e.CoffeeName, &e.GrindSize, &e.InputWeight,
		&e.OutputWeight, &e.TimeSeconds, &e.Rating, &e.Notes,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *DiaryEntryRepository) Create(ctx context.Context, e *domain.DiaryEntry) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO diary_entries
		(user_id, date_time, sieve_id, temperature, coffee_id, grind_size,
		 input_weight, output_weight, time_seconds, rating, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.UserID, e.DateTime, e.SieveID, e.Temperature, e.CoffeeID, e.GrindSize,
		e.InputWeight, e.OutputWeight, e.TimeSeconds, e.Rating, e.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DiaryEntryRepository) Update(ctx context.Context, e *domain.DiaryEntry) error {
	_, err := r.db.ExecContext(ctx, `UPDATE diary_entries SET
		date_time = ?, sieve_id = ?, temperature = ?, coffee_id = ?, grind_size = ?,
		input_weight = ?, output_weight = ?, time_seconds = ?, rating = ?, notes = ?
		WHERE id = ?`,
		e.DateTime, e.SieveID, e.Temperature, e.CoffeeID, e.GrindSize,
		e.InputWeight, e.OutputWeight, e.TimeSeconds, e.Rating, e.Notes, e.ID,
	)
	return err
}

func (r *DiaryEntryRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM diary_entries WHERE id = ?", id)
	return err
}

func buildWhere(f DiaryEntryFilter) (string, []interface{}) {
	clauses := []string{"de.user_id = ?"}
	args := []interface{}{f.UserID}

	if f.CoffeeID != nil {
		clauses = append(clauses, "de.coffee_id = ?")
		args = append(args, *f.CoffeeID)
	}
	if f.SieveID != nil {
		clauses = append(clauses, "de.sieve_id = ?")
		args = append(args, *f.SieveID)
	}
	if f.DateFrom != nil {
		clauses = append(clauses, "de.date_time >= ?")
		args = append(args, *f.DateFrom)
	}
	if f.DateTo != nil {
		clauses = append(clauses, "de.date_time <= ?")
		args = append(args, *f.DateTo)
	}
	if f.RatingMin != nil {
		clauses = append(clauses, "de.rating >= ?")
		args = append(args, *f.RatingMin)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}
