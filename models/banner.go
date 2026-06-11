package models

import "gorm.io/gorm"

type BannerStatus string

const (
	BannerStatusActive   BannerStatus = "active"
	BannerStatusInactive BannerStatus = "inactive"
)

type Banner struct {
	gorm.Model
	Title       string       `gorm:"type:varchar(150);not null"               json:"title"`
	Description *string      `gorm:"type:varchar(2000);default:null"          json:"description,omitempty"`
	ImageURL    string       `gorm:"type:varchar(500);not null"               json:"image_url"`
	LinkURL     *string      `gorm:"type:varchar(500);default:null"           json:"link_url,omitempty"`
	SortOrder   int          `gorm:"default:0"                               json:"sort_order"`
	Status      BannerStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedByID *uint        `gorm:"default:null"                             json:"created_by_id,omitempty"`
	CreatedBy   *User        `gorm:"foreignKey:CreatedByID"                   json:"created_by,omitempty"`
}

func (Banner) TableName() string { return "banners" }
