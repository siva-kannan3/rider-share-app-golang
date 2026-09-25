package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(r domain.TripRepository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)

}

func (svc *service) GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	osrmUrl := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=true&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	resp, err := http.Get(osrmUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes from OSRM: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routesResponse types.OsrmApiResponse

	if err := json.Unmarshal(body, &routesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse the response: %v", err)
	}

	return &routesResponse, nil
}
