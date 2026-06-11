package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name       string      `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	ImageURL   *string     `gorm:"type:varchar(500);default:null"         json:"image_url,omitempty"`
	SportClubs []SportClub `gorm:"many2many:sport_club_categories;"       json:"sport_clubs,omitempty"`
}

func (Category) TableName() string { return "categories" }
