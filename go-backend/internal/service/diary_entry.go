// Migrated from: DiaryEntryService.java
package service

import (
	"context"
	"database/sql"
	"math"

	"github.com/oglimmer/coffee-diary-backend/internal/domain"
	apperr "github.com/oglimmer/coffee-diary-backend/internal/errors"
	"github.com/oglimmer/coffee-diary-backend/internal/repository"
)

type DiaryEntryService struct {
	entryRepo  *repository.DiaryEntryRepository
	coffeeRepo *repository.CoffeeRepository
	sieveRepo  *repository.SieveRepository
}

func NewDiaryEntryService(
	entryRepo *repository.DiaryEntryRepository,
	coffeeRepo *repository.CoffeeRepository,
	sieveRepo *repository.SieveRepository,
) *DiaryEntryService {
	return &DiaryEntryService{
		entryRepo:  entryRepo,
		coffeeRepo: coffeeRepo,
		sieveRepo:  sieveRepo,
	}
}

func (s *DiaryEntryService) FindAll(ctx context.Context, filter repository.DiaryEntryFilter, page, size int, sortField, sortDir string) (*domain.PageResponse, error) {
	entries, total, err := s.entryRepo.FindAll(ctx, filter, page, size, sortField, sortDir)
	if err != nil {
		return nil, err
	}

	content := make([]domain.DiaryEntryResponse, 0, len(entries))
	for _, e := range entries {
		content = append(content, toResponse(e))
	}

	totalPages := 0
	if size > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(size)))
	}

	return &domain.PageResponse{
		Content:       content,
		TotalElements: total,
		TotalPages:    totalPages,
		Number:        page,
		Size:          size,
	}, nil
}

func (s *DiaryEntryService) FindByID(ctx context.Context, userID, entryID int64) (*domain.DiaryEntryResponse, error) {
	entry, err := s.entryRepo.FindByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, apperr.BadRequest("Diary entry not found")
	}
	if entry.UserID != userID {
		return nil, apperr.Forbidden("Not authorized to view this entry")
	}
	resp := toResponse(*entry)
	return &resp, nil
}

func (s *DiaryEntryService) Create(ctx context.Context, userID int64, req domain.DiaryEntryRequest) (*domain.DiaryEntryResponse, error) {
	temp := 93
	if req.Temperature != nil {
		temp = *req.Temperature
	}

	entry := &domain.DiaryEntry{
		UserID:      userID,
		DateTime:    req.DateTime.Time,
		Temperature: temp,
		GrindSize:   toNullFloat64(req.GrindSize),
		InputWeight: toNullFloat64(req.InputWeight),
		OutputWeight: toNullFloat64(req.OutputWeight),
		TimeSeconds: toNullInt64FromInt(req.TimeSeconds),
		Rating:      toNullInt64FromInt(req.Rating),
		Notes:       toNullString(req.Notes),
	}

	if req.SieveID != nil {
		sieve, err := s.sieveRepo.FindByID(ctx, *req.SieveID)
		if err != nil {
			return nil, err
		}
		if sieve == nil {
			return nil, apperr.BadRequest("Sieve not found")
		}
		if sieve.UserID != userID {
			return nil, apperr.BadRequest("Sieve does not belong to user")
		}
		entry.SieveID = sql.NullInt64{Int64: sieve.ID, Valid: true}
	}

	if req.CoffeeID != nil {
		coffee, err := s.coffeeRepo.FindByID(ctx, *req.CoffeeID)
		if err != nil {
			return nil, err
		}
		if coffee == nil {
			return nil, apperr.BadRequest("Coffee not found")
		}
		if coffee.UserID != userID {
			return nil, apperr.BadRequest("Coffee does not belong to user")
		}
		entry.CoffeeID = sql.NullInt64{Int64: coffee.ID, Valid: true}
	}

	id, err := s.entryRepo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Re-fetch to get joined sieve/coffee names
	created, err := s.entryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toResponse(*created)
	return &resp, nil
}

func (s *DiaryEntryService) Update(ctx context.Context, userID, entryID int64, req domain.DiaryEntryRequest) (*domain.DiaryEntryResponse, error) {
	entry, err := s.entryRepo.FindByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, apperr.BadRequest("Diary entry not found")
	}
	if entry.UserID != userID {
		return nil, apperr.Forbidden("Not authorized to update this entry")
	}

	temp := 93
	if req.Temperature != nil {
		temp = *req.Temperature
	}

	entry.DateTime = req.DateTime.Time
	entry.Temperature = temp
	entry.GrindSize = toNullFloat64(req.GrindSize)
	entry.InputWeight = toNullFloat64(req.InputWeight)
	entry.OutputWeight = toNullFloat64(req.OutputWeight)
	entry.TimeSeconds = toNullInt64FromInt(req.TimeSeconds)
	entry.Rating = toNullInt64FromInt(req.Rating)
	entry.Notes = toNullString(req.Notes)

	if req.SieveID != nil {
		sieve, err := s.sieveRepo.FindByID(ctx, *req.SieveID)
		if err != nil {
			return nil, err
		}
		if sieve == nil {
			return nil, apperr.BadRequest("Sieve not found")
		}
		if sieve.UserID != userID {
			return nil, apperr.BadRequest("Sieve does not belong to user")
		}
		entry.SieveID = sql.NullInt64{Int64: sieve.ID, Valid: true}
	} else {
		entry.SieveID = sql.NullInt64{}
	}

	if req.CoffeeID != nil {
		coffee, err := s.coffeeRepo.FindByID(ctx, *req.CoffeeID)
		if err != nil {
			return nil, err
		}
		if coffee == nil {
			return nil, apperr.BadRequest("Coffee not found")
		}
		if coffee.UserID != userID {
			return nil, apperr.BadRequest("Coffee does not belong to user")
		}
		entry.CoffeeID = sql.NullInt64{Int64: coffee.ID, Valid: true}
	} else {
		entry.CoffeeID = sql.NullInt64{}
	}

	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, err
	}

	updated, err := s.entryRepo.FindByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	resp := toResponse(*updated)
	return &resp, nil
}

func (s *DiaryEntryService) Delete(ctx context.Context, userID, entryID int64) error {
	entry, err := s.entryRepo.FindByID(ctx, entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return apperr.BadRequest("Diary entry not found")
	}
	if entry.UserID != userID {
		return apperr.Forbidden("Not authorized to delete this entry")
	}
	return s.entryRepo.Delete(ctx, entryID)
}

func toResponse(e domain.DiaryEntry) domain.DiaryEntryResponse {
	return domain.DiaryEntryResponse{
		ID:           e.ID,
		UserID:       e.UserID,
		DateTime:     domain.LocalDateTime{Time: e.DateTime},
		SieveID:      domain.NullInt64JSON{NullInt64: e.SieveID},
		SieveName:    domain.NullStringJSON{NullString: e.SieveName},
		Temperature:  e.Temperature,
		CoffeeID:     domain.NullInt64JSON{NullInt64: e.CoffeeID},
		CoffeeName:   domain.NullStringJSON{NullString: e.CoffeeName},
		GrindSize:    domain.NullFloat64JSON{NullFloat64: e.GrindSize},
		InputWeight:  domain.NullFloat64JSON{NullFloat64: e.InputWeight},
		OutputWeight: domain.NullFloat64JSON{NullFloat64: e.OutputWeight},
		TimeSeconds:  domain.NullInt64JSON{NullInt64: e.TimeSeconds},
		Rating:       domain.NullInt64JSON{NullInt64: e.Rating},
		Notes:        domain.NullStringJSON{NullString: e.Notes},
	}
}

func toNullFloat64(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}

func toNullInt64FromInt(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func toNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}
