package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	"sync"
)

type InMemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel

	mu sync.RWMutex
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *InMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.trips[trip.ID.Hex()] = trip

	return trip, nil
}

func (r *InMemRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.rideFares[fare.ID.Hex()] = fare

	return nil
}

func (r *InMemRepository) GetRideFareByID(ctx context.Context, fareId string) (*domain.RideFareModel, error) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	fare, exists := r.rideFares[fareId]
	if !exists {
		return nil, fmt.Errorf("Ride fare not found")
	}

	return fare, nil
}
