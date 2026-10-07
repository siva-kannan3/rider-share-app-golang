package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	internal_types "ride-sharing/services/trip-service/pkg/types"
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

func (svc *service) GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*internal_types.OsrmApiResponse, error) {
	osrmUrl := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?geometries=geojson",
		pickup.Longitude,
		pickup.Latitude,
		destination.Longitude,
		destination.Latitude,
	)

	log.Println(osrmUrl)
	resp, err := http.Get(osrmUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes from OSRM: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routesResponse internal_types.OsrmApiResponse

	if err := json.Unmarshal(body, &routesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse the response: %v", err)
	}

	return &routesResponse, nil
}

func (svc *service) EstimatePackagesPriceWithRoute(route *internal_types.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimateFares := make([]*domain.RideFareModel, len(baseFares))

	for i, baseFare := range baseFares {
		estimateFares[i] = estimateFareRoute(baseFare, route)
	}

	return estimateFares
}

func (svc *service) GenerateRideFares(ctx context.Context, fares []*domain.RideFareModel, userId string) ([]*domain.RideFareModel, error) {
	computedFares := make([]*domain.RideFareModel, len(fares))

	for i, fare := range fares {
		id := primitive.NewObjectID()
		cFare := domain.RideFareModel{
			UserID:            userId,
			ID:                id,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}

		if err := svc.repo.SaveRideFare(ctx, &cFare); err != nil {
			return nil, fmt.Errorf("Failed to save ride fare: %v", err)
		}
		computedFares[i] = &cFare
	}

	return computedFares, nil
}

func (svc *service) GetAndValidateFare(ctx context.Context, fareId string, userID string) (*domain.RideFareModel, error) {
	fare, err := svc.repo.GetRideFareByID(ctx, fareId)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip fare: %v", err)
	}

	// validating if user is matching with fetched fare
	if fare.UserID != userID {
		return nil, fmt.Errorf("Fare is not matching with userID")
	}

	return fare, nil
}

func estimateFareRoute(f *domain.RideFareModel, route *internal_types.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := internal_types.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationMins := route.Routes[0].Duration

	// distance fare
	distancePrice := distanceKm * pricingCfg.PricingPerUnitOfDistance

	// time fare
	timePrice := durationMins * pricingCfg.PricingPerMinute

	// total price
	totalPrice := distancePrice + timePrice + carPackagePrice

	return &domain.RideFareModel{
		PackageSlug:       f.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{

			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
		{

			PackageSlug:       "sedan",
			TotalPriceInCents: 400,
		},
		{

			PackageSlug:       "suv",
			TotalPriceInCents: 600,
		},
		{

			PackageSlug:       "van",
			TotalPriceInCents: 700,
		},
	}
}
