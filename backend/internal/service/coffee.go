// Migrated from: CoffeeService.java
package service

import (
	"context"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

type CoffeeService struct {
	coffeeRepo *repository.CoffeeRepository
}

func NewCoffeeService(coffeeRepo *repository.CoffeeRepository) *CoffeeService {
	return &CoffeeService{coffeeRepo: coffeeRepo}
}

func (s *CoffeeService) FindAllByUser(ctx context.Context, userID int64) ([]domain.CoffeeResponse, error) {
	coffees, err := s.coffeeRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CoffeeResponse, 0, len(coffees))
	for _, c := range coffees {
		result = append(result, domain.CoffeeResponse{ID: c.ID, Name: c.Name})
	}
	return result, nil
}

func (s *CoffeeService) Create(ctx context.Context, userID int64, req domain.CoffeeRequest) (*domain.CoffeeResponse, error) {
	if req.Name == "" {
		return nil, apperr.BadRequest("Name is required")
	}
	coffee, err := s.coffeeRepo.Create(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}
	return &domain.CoffeeResponse{ID: coffee.ID, Name: coffee.Name}, nil
}

func (s *CoffeeService) Delete(ctx context.Context, userID, coffeeID int64) error {
	coffee, err := s.coffeeRepo.FindByID(ctx, coffeeID)
	if err != nil {
		return err
	}
	if coffee == nil {
		return apperr.BadRequest("Coffee not found")
	}
	if coffee.UserID != userID {
		return apperr.Forbidden("Not authorized to delete this coffee")
	}
	return s.coffeeRepo.Delete(ctx, coffeeID)
}
