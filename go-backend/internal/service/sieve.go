// Migrated from: SieveService.java
package service

import (
	"context"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

type SieveService struct {
	sieveRepo *repository.SieveRepository
}

func NewSieveService(sieveRepo *repository.SieveRepository) *SieveService {
	return &SieveService{sieveRepo: sieveRepo}
}

func (s *SieveService) FindAllByUser(ctx context.Context, userID int64) ([]domain.SieveResponse, error) {
	sieves, err := s.sieveRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.SieveResponse, 0, len(sieves))
	for _, sv := range sieves {
		result = append(result, domain.SieveResponse{ID: sv.ID, Name: sv.Name})
	}
	return result, nil
}

func (s *SieveService) Create(ctx context.Context, userID int64, req domain.SieveRequest) (*domain.SieveResponse, error) {
	if req.Name == "" {
		return nil, apperr.BadRequest("Name is required")
	}
	sieve, err := s.sieveRepo.Create(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}
	return &domain.SieveResponse{ID: sieve.ID, Name: sieve.Name}, nil
}

func (s *SieveService) Delete(ctx context.Context, userID, sieveID int64) error {
	sieve, err := s.sieveRepo.FindByID(ctx, sieveID)
	if err != nil {
		return err
	}
	if sieve == nil {
		return apperr.BadRequest("Sieve not found")
	}
	if sieve.UserID != userID {
		return apperr.Forbidden("Not authorized to delete this sieve")
	}
	return s.sieveRepo.Delete(ctx, sieveID)
}
