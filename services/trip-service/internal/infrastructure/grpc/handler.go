package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer

	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, svc domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: svc,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()
	userID := req.GetUserID()

	pickupCoords := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}

	destinationCoords := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}

	resp, err := h.service.GetRoute(ctx, pickupCoords, destinationCoords)
	log.Println(resp, err)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get Route: %v", err)
	}

	estimatedFares := h.service.EstimatePackagesPriceWithRoute(resp)

	rideFares, err := h.service.GenerateRideFares(ctx, estimatedFares, userID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to generate ride fare: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     resp.ToProto(),
		RideFares: domain.ToRideFaresProto(rideFares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	rideFareId := req.GetRideFareID()
	userID := req.GetUserID()
	// Get fare
	fare, err := h.service.GetAndValidateFare(ctx, rideFareId, userID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get ride fare: %v", err)
	}

	// create trip
	_, err = h.service.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: "123",
	}, nil
}
