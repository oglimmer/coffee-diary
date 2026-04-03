package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
)

func TestToResponse(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	entry := domain.DiaryEntry{
		ID:          1,
		UserID:      2,
		DateTime:    now,
		SieveID:     sql.NullInt64{Int64: 5, Valid: true},
		SieveName:   sql.NullString{String: "IMS", Valid: true},
		Temperature: 93,
		CoffeeID:    sql.NullInt64{Int64: 3, Valid: true},
		CoffeeName:  sql.NullString{String: "Ethiopian", Valid: true},
		GrindSize:   sql.NullFloat64{Float64: 5.0, Valid: true},
		InputWeight: sql.NullFloat64{Float64: 18.0, Valid: true},
		OutputWeight: sql.NullFloat64{Float64: 36.0, Valid: true},
		TimeSeconds: sql.NullInt64{Int64: 25, Valid: true},
		Rating:      sql.NullInt64{Int64: 4, Valid: true},
		Notes:       sql.NullString{String: "Great shot", Valid: true},
	}

	resp := toResponse(entry)

	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, int64(2), resp.UserID)
	assert.Equal(t, now, resp.DateTime.Time)
	assert.True(t, resp.SieveID.Valid)
	assert.Equal(t, int64(5), resp.SieveID.Int64)
	assert.Equal(t, "IMS", resp.SieveName.String)
	assert.Equal(t, 93, resp.Temperature)
	assert.True(t, resp.CoffeeID.Valid)
	assert.Equal(t, int64(3), resp.CoffeeID.Int64)
	assert.Equal(t, "Ethiopian", resp.CoffeeName.String)
	assert.Equal(t, 5.0, resp.GrindSize.Float64)
	assert.Equal(t, 18.0, resp.InputWeight.Float64)
	assert.Equal(t, 36.0, resp.OutputWeight.Float64)
	assert.Equal(t, int64(25), resp.TimeSeconds.Int64)
	assert.Equal(t, int64(4), resp.Rating.Int64)
	assert.Equal(t, "Great shot", resp.Notes.String)
}

func TestToResponse_NullFields(t *testing.T) {
	entry := domain.DiaryEntry{
		ID:          1,
		UserID:      2,
		DateTime:    time.Now(),
		Temperature: 93,
	}

	resp := toResponse(entry)

	assert.False(t, resp.SieveID.Valid)
	assert.False(t, resp.SieveName.Valid)
	assert.False(t, resp.CoffeeID.Valid)
	assert.False(t, resp.CoffeeName.Valid)
	assert.False(t, resp.GrindSize.Valid)
	assert.False(t, resp.InputWeight.Valid)
	assert.False(t, resp.OutputWeight.Valid)
	assert.False(t, resp.TimeSeconds.Valid)
	assert.False(t, resp.Rating.Valid)
	assert.False(t, resp.Notes.Valid)
}

func TestToNullFloat64(t *testing.T) {
	v := 5.5
	result := toNullFloat64(&v)
	assert.True(t, result.Valid)
	assert.Equal(t, 5.5, result.Float64)

	result = toNullFloat64(nil)
	assert.False(t, result.Valid)
}

func TestToNullInt64FromInt(t *testing.T) {
	v := 42
	result := toNullInt64FromInt(&v)
	assert.True(t, result.Valid)
	assert.Equal(t, int64(42), result.Int64)

	result = toNullInt64FromInt(nil)
	assert.False(t, result.Valid)
}

func TestToNullString(t *testing.T) {
	v := "hello"
	result := toNullString(&v)
	assert.True(t, result.Valid)
	assert.Equal(t, "hello", result.String)

	result = toNullString(nil)
	assert.False(t, result.Valid)
}
