package domain

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDateTime_MarshalJSON(t *testing.T) {
	dt := LocalDateTime{Time: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)}
	b, err := json.Marshal(dt)
	require.NoError(t, err)
	assert.Equal(t, `"2024-01-15T10:30:00"`, string(b))
}

func TestLocalDateTime_UnmarshalJSON(t *testing.T) {
	var dt LocalDateTime
	err := json.Unmarshal([]byte(`"2024-01-15T10:30:00"`), &dt)
	require.NoError(t, err)
	assert.Equal(t, 2024, dt.Time.Year())
	assert.Equal(t, time.January, dt.Time.Month())
	assert.Equal(t, 15, dt.Time.Day())
	assert.Equal(t, 10, dt.Time.Hour())
	assert.Equal(t, 30, dt.Time.Minute())
}

func TestNullFloat64JSON_Marshal(t *testing.T) {
	valid := NullFloat64JSON{sql.NullFloat64{Float64: 18.5, Valid: true}}
	b, _ := json.Marshal(valid)
	assert.Equal(t, "18.5", string(b))

	null := NullFloat64JSON{sql.NullFloat64{}}
	b, _ = json.Marshal(null)
	assert.Equal(t, "null", string(b))
}

func TestNullInt64JSON_Marshal(t *testing.T) {
	valid := NullInt64JSON{sql.NullInt64{Int64: 42, Valid: true}}
	b, _ := json.Marshal(valid)
	assert.Equal(t, "42", string(b))

	null := NullInt64JSON{sql.NullInt64{}}
	b, _ = json.Marshal(null)
	assert.Equal(t, "null", string(b))
}

func TestNullStringJSON_Marshal(t *testing.T) {
	valid := NullStringJSON{sql.NullString{String: "hello", Valid: true}}
	b, _ := json.Marshal(valid)
	assert.Equal(t, `"hello"`, string(b))

	null := NullStringJSON{sql.NullString{}}
	b, _ = json.Marshal(null)
	assert.Equal(t, "null", string(b))
}

func TestDiaryEntryResponse_JSON(t *testing.T) {
	resp := DiaryEntryResponse{
		ID:           1,
		UserID:       2,
		DateTime:     LocalDateTime{Time: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)},
		SieveID:      NullInt64JSON{sql.NullInt64{Int64: 5, Valid: true}},
		SieveName:    NullStringJSON{sql.NullString{String: "IMS", Valid: true}},
		Temperature:  93,
		CoffeeID:     NullInt64JSON{sql.NullInt64{}},
		CoffeeName:   NullStringJSON{sql.NullString{}},
		GrindSize:    NullFloat64JSON{sql.NullFloat64{Float64: 5.0, Valid: true}},
		InputWeight:  NullFloat64JSON{sql.NullFloat64{Float64: 18.0, Valid: true}},
		OutputWeight: NullFloat64JSON{sql.NullFloat64{Float64: 36.0, Valid: true}},
		TimeSeconds:  NullInt64JSON{sql.NullInt64{Int64: 25, Valid: true}},
		Rating:       NullInt64JSON{sql.NullInt64{Int64: 4, Valid: true}},
		Notes:        NullStringJSON{sql.NullString{}},
	}

	b, err := json.Marshal(resp)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	assert.Equal(t, float64(1), m["id"])
	assert.Equal(t, float64(2), m["userId"])
	assert.Equal(t, "2024-01-15T10:00:00", m["dateTime"])
	assert.Equal(t, float64(5), m["sieveId"])
	assert.Equal(t, "IMS", m["sieveName"])
	assert.Equal(t, float64(93), m["temperature"])
	assert.Nil(t, m["coffeeId"])
	assert.Nil(t, m["coffeeName"])
	assert.Equal(t, float64(5), m["grindSize"])
	assert.Equal(t, float64(18), m["inputWeight"])
	assert.Equal(t, float64(36), m["outputWeight"])
	assert.Equal(t, float64(25), m["timeSeconds"])
	assert.Equal(t, float64(4), m["rating"])
	assert.Nil(t, m["notes"])
}

func TestPageResponse_JSON(t *testing.T) {
	resp := PageResponse{
		Content:       []CoffeeResponse{{ID: 1, Name: "Ethiopian"}},
		TotalElements: 1,
		TotalPages:    1,
		Number:        0,
		Size:          20,
	}

	b, err := json.Marshal(resp)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	assert.Equal(t, float64(1), m["totalElements"])
	assert.Equal(t, float64(1), m["totalPages"])
	assert.Equal(t, float64(0), m["number"])
	assert.Equal(t, float64(20), m["size"])
}
