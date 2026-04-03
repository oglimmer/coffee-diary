package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
)

func TestSieveService_Create_EmptyName(t *testing.T) {
	svc := &SieveService{sieveRepo: nil}
	_, err := svc.Create(context.Background(), 1, domain.SieveRequest{Name: ""})
	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, 400, appErr.Status)
	assert.Contains(t, appErr.Message, "Name is required")
}
