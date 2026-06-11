package models

import "gorm.io/gorm"

type SportClub struct {
	gorm.Model
	Name          string     `gorm:"type:varchar(150);not null"               json:"name"`
	Latitude      float64    `gorm:"type:decimal(10,8);not null"              json:"lat"`
	Longitude     float64    `gorm:"type:decimal(11,8);not null"              json:"lng"`
	Location      string     `gorm:"type:varchar(255);not null"               json:"location"`
	IsOpen        bool       `gorm:"default:false"                            json:"is_open"`
	OpenTime      *string    `gorm:"type:varchar(10);default:null"            json:"open_time,omitempty"`
	CloseTime     *string    `gorm:"type:varchar(10);default:null"            json:"close_time,omitempty"`
	Description   *string    `gorm:"type:varchar(2000);default:null"          json:"description,omitempty"`
	ImageURLs     string     `gorm:"type:varchar(2000);not null;default:'[]'" json:"image_urls"`
	FavoriteCount int        `gorm:"default:0"                                json:"favorite_count"`
	CreatedByID   *uint      `gorm:"default:null"                             json:"created_by_id,omitempty"`
	CreatedBy     *User      `gorm:"foreignKey:CreatedByID"                   json:"created_by,omitempty"`
	Categories    []Category `gorm:"many2many:sport_club_categories;"         json:"categories,omitempty"`
	Slots         []Slot     `gorm:"foreignKey:SportClubID"                   json:"slots,omitempty"`
}

func (SportClub) TableName() string { return "sport_clubs" }
