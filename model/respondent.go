package model

import (
	"time"

	os "os"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"gorm.io/gorm"
)

type Respondent struct {
	// -
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Active    bool      `gorm:"default:false;not null" json:"active"`
	Published bool      `gorm:"default:false;not null" json:"published"`

	//PRIMARY VEHICLE USED BY THIS RESPONDENT
	VehicleID *uuid.UUID `gorm:"unique" json:"vehicleId,omitempty"`
	Vehicle   *Vehicle   `gorm:"foreignKey:VehicleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vehicle,omitempty"`

	BillingAmount *uint64 `gorm:"" json:"billingAmount"`

	AccessBalance      *RespondentAccessProductBalance      `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accessBalance,omitempty"`
	AccessSubscription *RespondentAccessProductSubscription `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accessSubscription,omitempty"`

	// //The last location entry posted by the session
	// CurrentLocationEntryID *uuid.UUID                      `gorm:"unique" json:"currentLocationEntryId,omitempty"`
	// CurrentLocationEntry   *RespondentSessionLocationEntry `gorm:"foreignKey:CurrentLocationEntryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"currentLocationEntry,omitempty"`

	// >> DOCUMENTS AND VERIFICATIONS
	Ssn       *string           `gorm: "type:varchar(32);unique" json:"ssn,omitempty"`
	Documents *[]*ImageDocument `gorm:"foreignKey:ResponderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"documents,omitempty"`
	// DOCUMENTS

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:";unique" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	//COMPANY OF THE REPSONDANT
	CompanyID *uuid.UUID `gorm:"" json:"companyId,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company,omitempty"`
	// PLACE THE RESPONDENT CAN OPERATE IN
	PlaceID *uuid.UUID `gorm:"" json:"placeId,omitempty"`
	Place   *Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"place,omitempty"`
	// TIMESTAMPS
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (Respondent *Respondent) BeforeCreate(tx *gorm.DB) (err error) {
	Respondent.ID = uuid.New()
	Respondent.CreatedAt = time.Now()
	Respondent.UpdatedAt = time.Now()
	return nil
}

func (Respondent *Respondent) Independent() bool {
	return Respondent.CompanyID == nil
}

func (responder *Respondent) ClearDocuments(types ...string) (err error) {
	// if len(types) == 0 {
	// 	responder.Documents = nil
	// 	return nil
	// }
	config, err := initializers.LoadConfig()
	if err != nil {
		return err
	}

	var fullRespondent Respondent
	if err = DB.Where("id = ?", responder.ID).Preload("Documents").First(&fullRespondent).Error; err != nil {
		return err
	}

	var tLength = len(types)
	var newDocuments []*ImageDocument
	for _, doc := range *fullRespondent.Documents {
		docType := doc.DocType
		if tLength == 0 || (docType != nil && contains(types, *docType)) {
			//REMOVE THIS DOCUMENT
			nativePath := config.ResolveNativePath(doc.FilePath)
			// Delete File
			if err = os.Remove(nativePath); err != nil {
				return err
			}
			if err = DB.Delete(doc).Error; err != nil {
				return err
			}

		} else {
			newDocuments = append(newDocuments, doc)
		}
	}
	responder.Documents = &newDocuments
	return nil
}

func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
