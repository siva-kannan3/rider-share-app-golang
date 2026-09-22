package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()
	inmemRepo := repository.NewInMemRepository()

	svc := service.NewService(inmemRepo)

	fare := &domain.RideFareModel{
		ID:                primitive.NewObjectID(),
		UserID:            "223",
		PackageSlug:       "sedan",
		TotalPriceInCents: 240,
		ExpiresAt:         time.Now(),
	}

	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}

	log.Println(t)

}
