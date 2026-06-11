package models

import "gorm.io/gorm"

type Slot struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(150);not null"      json:"name"`
	ImageURL    *string   `gorm:"type:varchar(500);default:null"  json:"image_url,omitempty"`
	Description *string   `gorm:"type:varchar(2000);default:null" json:"description,omitempty"`
	Price       float64   `gorm:"type:decimal(10,2);default:0"    json:"price"`
	Capacity    int       `gorm:"default:1"                       json:"capacity"`
	IsAvailable bool      `gorm:"default:true"                    json:"is_available"`
	SportClubID uint      `gorm:"not null;index"                  json:"sport_club_id"`
	SportClub   SportClub `gorm:"foreignKey:SportClubID"          json:"sport_club,omitempty"`
	CategoryID  *uint     `gorm:"default:null;index"              json:"category_id,omitempty"`
	Category    *Category `gorm:"foreignKey:CategoryID"           json:"category,omitempty"`
	CreatedByID *uint     `gorm:"default:null"                    json:"created_by_id,omitempty"`
	CreatedBy   *User     `gorm:"foreignKey:CreatedByID"          json:"created_by,omitempty"`
}

func (Slot) TableName() string { return "slots" }
