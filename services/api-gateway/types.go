package main

import "ride-sharing/shared/types"

type tripPreviewRequest struct {
	UserID      string           `json:"UserID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}
