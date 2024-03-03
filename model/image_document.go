package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"gorm.io/gorm"
)

type ImageDocument struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	DocType      *string   `gorm:"type:varchar(32);" json:"docType,omitempty"`
	FilePath     string    `gorm:"type:varchar(1000);not null" json:"filePath,omitempty"`
	Size         int64     `gorm:"not null" json:"size,omitempty"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"originalName,omitempty"`
	Extension    string    `gorm:"type:varchar(10);not null" json:"extension,omitempty"`
	Verified     *bool     `gorm:"default:false;not null" json:"verified"`
	//
	VehicleID *uuid.UUID `gorm:"type:char(36);index:idx_image_document_vehicle_id,LENGTH(36);" json:"vehicleId,omitempty"`
	Vehicle   *Vehicle   `gorm:"foreignKey:VehicleID;references:ID" json:"vehicle,omitempty"`
	//
	ResponderID *uuid.UUID  `gorm:"type:char(36);index:idx_image_document_responder_id,LENGTH(36);" json:"responderId,omitempty"`
	Responder   *Respondent `gorm:"foreignKey:ResponderID;references:ID" json:"responder,omitempty"`
	//

	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (document *ImageDocument) BeforeCreate(tx *gorm.DB) (err error) {

	if document.ResponderID == nil && document.VehicleID == nil {
		return errors.New("ImageDocument must have a vehicle or a responder")
	}
	if document.ResponderID != nil && document.VehicleID != nil {
		return errors.New("ImageDocument must have either vehicle or a responder but not both")
	}

	document.ID = uuid.New()
	document.CreatedAt = time.Now()
	document.UpdatedAt = time.Now()
	return nil
}

func (document *ImageDocument) PublicPath() string {
	config, err := initializers.LoadConfig()
	if err != nil {
		// return err
		message := "Error loading config"
		panic(message)
	}
	return config.ResolvePublicPath(document.FilePath)
}

func ResolveDocumentSlicePublicPaths(documents *[]*ImageDocument) {
	config, err := initializers.LoadConfig()
	if err != nil {
		// return err
		message := "Error loading config"
		panic(message)
	}
	if documents != nil {
		for _, doc := range *documents {
			doc.FilePath = config.ResolvePublicPath(doc.FilePath)
		}
	}
}

// for _, doc := range documents {
// 	doc.FilePath = config.ResolvePublicPath(doc.FilePath)
// }
