package domain

import (
	"context"
	internal_types "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
	GetRideFareByID(ctx context.Context, fareId string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*internal_types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *internal_types.OsrmApiResponse) []*RideFareModel
	GenerateRideFares(ctx context.Context, fares []*RideFareModel, userID string, route *internal_types.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareId string, userID string) (*RideFareModel, error)
}
