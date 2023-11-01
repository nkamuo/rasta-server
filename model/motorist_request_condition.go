package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MotoristRequestSituation struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//
	Title     string `gorm:"not null;type:varchar(128);" json:"title,omitempty"`
	SubTitlte string `gorm:"type:varchar(255);" json:"subtitle,omitempty"`
	//
	Code string `gorm:"unique;type:varchar(64);" json:"code,omitempty"`
	//
	Description string `gorm:"" json:"description,omitempty"`
	Note        string `gorm:"" json:"note,omitempty"`
	//
	// Active      *bool     `gorm:"default:false;not null" json:"active"`
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (situation *MotoristRequestSituation) BeforeCreate(tx *gorm.DB) (err error) {
	situation.ID = uuid.New()

	situation.CreatedAt = time.Now()
	situation.UpdatedAt = time.Now()
	return nil
}

var defaultConfitions = []MotoristRequestSituation{
	{
		// ID:          uuid.New(),
		Code:        "needs_to_ride_with_driver_on_the_front_seet",
		Title:       "Needs to ride with driver",
		SubTitlte:   "Needs to ride with driver on the front seat.",
		Description: "In the case of towing request, the motorist needs to sit on the driver with the passenger's seat",
		Note:        "The motorist is in a dark or low-light place",
	},
	{
		// ID:          uuid.New(),
		Code:        "in_a_dark_place",
		Title:       "In a Low light area",
		SubTitlte:   "In a Dark or low-light place",
		Description: "You are in a dark place",
		Note:        "The motorist is in a dark or low-light place",
	},
	{
		// ID:          uuid.New(),
		Code:        "on_high_way",
		Title:       "On the highway",
		SubTitlte:   "Currently On the high way or express",
		Description: "You are on the high way",
		Note:        "The motorist is on the highway",
	},
	{
		// ID:          uuid.New(),
		Code:        "on_access_road_or_street",
		Title:       "On access Road",
		SubTitlte:   "In the city street or access road",
		Description: "Your current location is int the street",
		Note:        "The motorist is on access road/street",
	},
	{
		// ID:          uuid.New(),
		Code:        "on_private_business_parking_lot",
		Title:       "On a parking lot",
		SubTitlte:   "On a private business parking lot",
		Description: "Your current location is int the street",
		Note:        "Currently in a parking lot",
	},
	{
		// ID:          uuid.New(),
		Code:        "vehicle_blocking_or_near_traffic_lane",
		Title:       "Blocking Traffic",
		SubTitlte:   "Vehicle Blocking or near Traffic Lane",
		Description: "Vehicle Blocking or near Traffic Lane",
		Note:        "Vehicle Blocking or near Traffic Lane",
	},
	{
		// ID:          uuid.New(),
		Code:        "vehicle_have_been_involved_in_an_accident",
		Title:       "Involved in an accident",
		SubTitlte:   "Vehicle have been invlolved in an accident",
		Description: "Vehicle have been invlolved in an accident",
		Note:        "Vehicle have been invlolved in an accident",
	},
}

func GetDefaultMotoristRequestSituations() (situations []MotoristRequestSituation) {
	return defaultConfitions
}

func MigrateMotoristSituations(db *gorm.DB) (err error) {
	// situationRepo := repository.GetMotoristRequestSituationRepository()

	situations := GetDefaultMotoristRequestSituations()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, situation := range situations {
			// tx.Model(&situation)
			if Situation, err := getByCode(tx, situation.Code); err != nil {
				if err.Error() == "record not found" {
					if err := tx.Save(&situation).Error; err != nil {
						return err
					}

				} else {
					return err
				}
			} else {
				Situation.Code = situation.Code
				Situation.Title = situation.Title
				Situation.SubTitlte = situation.SubTitlte
				Situation.Note = situation.Note
				Situation.Description = situation.Description

				if err := tx.Save(Situation).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func getByCode(db *gorm.DB, code string) (motoristRequestSituation *MotoristRequestSituation, err error) {
	if err = db.Where("code = ?", code).First(&motoristRequestSituation).Error; err != nil {
		return nil, err
	}
	return motoristRequestSituation, nil
}
