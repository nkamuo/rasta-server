package dto

import "github.com/google/uuid"

type LocationInput struct {
	ID          uuid.UUID                 `json:"id,omitempty" binding:"required"`
	Name        string                    `json:"name,omitempty" binding:"required"`
	Street      string                    `json:"longName,omitempty" binding:"required"`
	Address     string                    `json:"shortName,omitempty" binding:"required"`
	Coordinates *LocationCoordinatesInput `json:"coords" binding:"required"`
	GoogleID    *string                   `json:"googleId,omitempty" binding:"required"`
	Description *string                   `json:"description,omitempty" binding:""`
}

type LocationCoordinatesInput struct {
	Latitude  float32  `json:"latitude,omitempty" binding:"required"`
	Longitude float32  `json:"longitude,omitempty" binding:"required"`
	Altitude  float32  `json:"altitude,omitempty"`
	Accuracy  *float32 `json:"accuracy,omitempty"`
	Heading   *float32 `json:"heading,omitempty"`
	Speed     *float32 `json:"speed,omitempty"`
}

type TransitLocationRequestInput struct {
	Origin      string `json:"origin,omitempty;" form:"origin" binding:"required"`
	Destination string `json:"destination,omitempty;" form:"destination" binding:"required"`
}

type DistanceMatrixRequestInput struct {
	TransitLocationRequestInput
}
