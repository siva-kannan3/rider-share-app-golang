package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type tripPreviewRequest struct {
	UserID      string           `json:"UserID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (tr *tripPreviewRequest) toProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: tr.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  tr.Pickup.Latitude,
			Longitude: tr.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  tr.Destination.Latitude,
			Longitude: tr.Destination.Longitude,
		},
	}
}
