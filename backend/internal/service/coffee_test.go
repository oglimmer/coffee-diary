package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

// --- mock ---

type mockCoffeeDB struct {
	coffees map[int64]*domain.Coffee
	nextID  int64
}

func newMockCoffeeDB() *mockCoffeeDB {
	return &mockCoffeeDB{coffees: make(map[int64]*domain.Coffee), nextID: 1}
}

func setupCoffeeService(db *mockCoffeeDB) (*CoffeeService, *repository.CoffeeRepository) {
	// We can't easily mock the repository since it takes *sql.DB.
	// Instead we test at the handler/integration level for DB-dependent code.
	// These tests validate the service logic using the real service with a nil repo
	// where we can test validation paths.
	return nil, nil // placeholder — see below for validation tests
}

func TestCoffeeService_Create_EmptyName(t *testing.T) {
	svc := &CoffeeService{coffeeRepo: nil} // repo not called for validation failure
	_, err := svc.Create(context.Background(), 1, domain.CoffeeRequest{Name: ""})
	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, 400, appErr.Status)
	assert.Contains(t, appErr.Message, "Name is required")
}

func TestCoffeeService_Delete_NotFound(t *testing.T) {
	// With a nil repo this would panic, so we skip DB-dependent tests here.
	// Full coverage comes from integration tests.
	t.Skip("requires database — covered by integration tests")
}
