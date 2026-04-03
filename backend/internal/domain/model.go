// Migrated from: entity/*.java + dto/*.java
package domain

import (
	"database/sql"
	"encoding/json"
	"time"
)

// LocalDateTime matches Spring Boot's LocalDateTime JSON format (no timezone).
type LocalDateTime struct {
	time.Time
}

const localDateTimeFormat = "2006-01-02T15:04:05"

func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(localDateTimeFormat))
}

func (t *LocalDateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.Parse(localDateTimeFormat, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// NullFloat64JSON serializes sql.NullFloat64 as null or number.
type NullFloat64JSON struct {
	sql.NullFloat64
}

func (n NullFloat64JSON) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}

func (n *NullFloat64JSON) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(b, &n.Float64)
}

// NullInt64JSON serializes sql.NullInt64 as null or number.
type NullInt64JSON struct {
	sql.NullInt64
}

func (n NullInt64JSON) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}

func (n *NullInt64JSON) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(b, &n.Int64)
}

// NullStringJSON serializes sql.NullString as null or string.
type NullStringJSON struct {
	sql.NullString
}

func (n NullStringJSON) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.String)
}

func (n *NullStringJSON) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(b, &n.String)
}

// --- Entities ---

type User struct {
	ID        int64
	Username  string
	OIDCSub   string
	CreatedAt time.Time
}

type Coffee struct {
	ID     int64
	Name   string
	UserID int64
}

type Sieve struct {
	ID     int64
	Name   string
	UserID int64
}

type DiaryEntry struct {
	ID           int64
	UserID       int64
	DateTime     time.Time
	SieveID      sql.NullInt64
	SieveName    sql.NullString
	Temperature  int
	CoffeeID     sql.NullInt64
	CoffeeName   sql.NullString
	GrindSize    sql.NullFloat64
	InputWeight  sql.NullFloat64
	OutputWeight sql.NullFloat64
	TimeSeconds  sql.NullInt64
	Rating       sql.NullInt64
	Notes        sql.NullString
}

// --- Request DTOs ---

type CoffeeRequest struct {
	Name string `json:"name"`
}

type SieveRequest struct {
	Name string `json:"name"`
}

type DiaryEntryRequest struct {
	DateTime    LocalDateTime  `json:"dateTime"`
	SieveID     *int64         `json:"sieveId"`
	Temperature *int           `json:"temperature"`
	CoffeeID    *int64         `json:"coffeeId"`
	GrindSize   *float64       `json:"grindSize"`
	InputWeight *float64       `json:"inputWeight"`
	OutputWeight *float64      `json:"outputWeight"`
	TimeSeconds *int           `json:"timeSeconds"`
	Rating      *int           `json:"rating"`
	Notes       *string        `json:"notes"`
}

// --- Response DTOs ---

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type CoffeeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SieveResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DiaryEntryResponse struct {
	ID           int64           `json:"id"`
	UserID       int64           `json:"userId"`
	DateTime     LocalDateTime   `json:"dateTime"`
	SieveID      NullInt64JSON   `json:"sieveId"`
	SieveName    NullStringJSON  `json:"sieveName"`
	Temperature  int             `json:"temperature"`
	CoffeeID     NullInt64JSON   `json:"coffeeId"`
	CoffeeName   NullStringJSON  `json:"coffeeName"`
	GrindSize    NullFloat64JSON `json:"grindSize"`
	InputWeight  NullFloat64JSON `json:"inputWeight"`
	OutputWeight NullFloat64JSON `json:"outputWeight"`
	TimeSeconds  NullInt64JSON   `json:"timeSeconds"`
	Rating       NullInt64JSON   `json:"rating"`
	Notes        NullStringJSON  `json:"notes"`
}

type PageResponse struct {
	Content       interface{} `json:"content"`
	TotalElements int64       `json:"totalElements"`
	TotalPages    int         `json:"totalPages"`
	Number        int         `json:"number"`
	Size          int         `json:"size"`
}
