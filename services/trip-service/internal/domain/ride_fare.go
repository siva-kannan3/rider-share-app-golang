package domain

import (
	internal_types "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // e.g. sedan, luxury, auto
	TotalPriceInCents float64
	Route             *internal_types.OsrmApiResponse
	// ExpiresAt         time.Time
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	result := make([]*pb.RideFare, len(fares))

	for i, fare := range fares {
		result[i] = fare.ToProto()
	}

	return result
}
