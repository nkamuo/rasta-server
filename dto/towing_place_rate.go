package dto

import "github.com/google/uuid"

type TowingPlaceRateCreationInput struct {
	Rate        uint64  `json:"rate" binding:"required"`
	MinDistance *uint64 `json:"minDistance" binding:"required,gte=0,ltfield=MaxDistance"`
	MaxDistance uint64  `json:"maxDistance" binding:"required,gte=0,gtfield=MinDistance"`
	//
	Sequence uint64 `json:"sequence" binding:"gte=0"`
	//
	PlaceID uuid.UUID `json:"placeId,omitempty" binding:"required"`
	//
	Active      bool    `json:"active" binding:""`
	Description *string `json:"description,omitempty" binding:""`
}

type TowingPlaceRateUpdateInput struct {
	Rate        *uint64 `json:"rate" binding:"required"`
	MinDistance *uint64 `json:"minDistance" binding:"required"`
	MaxDistance *uint64 `json:"maxDistance" binding:"required"`
	//
	Sequence *uint64 `json:"sequence" binding:""`
	//
	Active      *bool   `json:"active" binding:""`
	Description *string `json:"description,omitempty" binding:""`
}

type TowingPlaceRateRequestQuery struct {
	Origin      string `json:"origin,omitempty" form:"origin" binding:"required"`
	Destination string `json:"destination,omitempty" form:"destination" binding:"required"`
}

type TowingPlaceRateDistanceRequestQuery struct {
	Distance uint64 `json:"distance" form:"distance" binding:"gte=0"`
	PlaceID  string `json:"placeId,omitempty" form:"place_id" binding:"required"`
}
