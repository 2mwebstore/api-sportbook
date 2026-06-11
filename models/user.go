package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleClient  Role = "client"
	RolePartner Role = "partner"
	RoleAdmin   Role = "admin"
)

type User struct {
	gorm.Model
	FullName  string  `gorm:"type:varchar(100);not null"                  json:"full_name"`
	Email     *string `gorm:"type:varchar(100);uniqueIndex;default:null"  json:"email,omitempty"`
	Phone     *string `gorm:"type:varchar(20);uniqueIndex;default:null"   json:"phone,omitempty"`
	Password  string  `gorm:"type:varchar(255);not null"                  json:"-"`
	AvatarURL *string `gorm:"type:varchar(500);default:null"              json:"avatar_url,omitempty"`
	Role      Role    `gorm:"type:varchar(20);not null;default:'client'"  json:"role"`
	// Location fields (client)
	Latitude    *float64   `gorm:"type:decimal(10,8);default:null"             json:"lat,omitempty"`
	Longitude   *float64   `gorm:"type:decimal(11,8);default:null"             json:"lng,omitempty"`
	Location    *string    `gorm:"type:varchar(255);default:null"              json:"location,omitempty"`
	IsVerified  bool       `gorm:"default:false"                               json:"is_verified"`
	IsActive    bool       `gorm:"default:true"                                json:"is_active"`
	LastLoginAt *time.Time `                                                   json:"last_login_at,omitempty"`
	// Favorites
	Favorites []SportClub `gorm:"many2many:user_favorites;"                   json:"favorites,omitempty"`
}

func (User) TableName() string { return "users" }
