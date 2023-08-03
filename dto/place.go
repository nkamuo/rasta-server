package dto

import "github.com/nkamuo/rasta-server/model"

type PlaceCreationInput struct {
	Code        string                 `json:"code,omitempty" binding:"required"`
	Name        string                 `json:"name,omitempty" binding:"required"`
	ShortName   string                 `json:"shortName,omitempty"`
	LongName    string                 `json:"longName,omitempty"`
	GoogleID    *string                `json:"goggleId,omitempty"`
	Coordinates model.PlaceCoordinates `json:"coords,omitempty" binding:"required"`
	Description string                 `json:"description,omitempty"`
	Category    model.PlaceCategory    `json:"category,omitempty" binding:"required"`
	Active      bool                   `json:"active"`
	// Level  int16 `json:"level" binding:"exists"`	//exists will allow zero values as valid unlike required which will fail
}

type PlaceUpdateInput struct {
	Code        *string                 `json:"code,omitempty" binding:""`
	Name        *string                 `json:"name,omitempty"`
	ShortName   *string                 ` json:"shortName,omitempty"`
	LongName    *string                 `json:"longName,omitempty"`
	Description *string                 `json:"description,omitempty"`
	GoogleID    *string                 `json:"goggleId,omitempty"`
	Coordinates *model.PlaceCoordinates `json:"coords,omitempty" binding:""`
	Category    *model.PlaceCategory    `json:"category,omitempty"`
	Active      *bool                   `json:"active"`
}
