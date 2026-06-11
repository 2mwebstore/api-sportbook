package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(255);default:''"           json:"description"`
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Role         string     `gorm:"type:varchar(20);not null;index" json:"role"`
	PermissionID uint       `gorm:"not null"                        json:"permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID"         json:"permission,omitempty"`
}

func (RolePermission) TableName() string { return "role_permissions" }

const (
	PermSportClubRead   = "sport_club:read"
	PermSportClubCreate = "sport_club:create"
	PermSportClubUpdate = "sport_club:update"
	PermSportClubDelete = "sport_club:delete"

	PermCategoryRead   = "category:read"
	PermCategoryCreate = "category:create"
	PermCategoryUpdate = "category:update"
	PermCategoryDelete = "category:delete"

	PermBannerRead   = "banner:read"
	PermBannerCreate = "banner:create"
	PermBannerUpdate = "banner:update"
	PermBannerDelete = "banner:delete"

	PermSlotRead   = "slot:read"
	PermSlotCreate = "slot:create"
	PermSlotUpdate = "slot:update"
	PermSlotDelete = "slot:delete"

	PermUserReadSelf   = "user:read_self"
	PermUserUpdateSelf = "user:update_self"
)

var AllPermissions = []Permission{
	{Name: PermSportClubRead, Description: "Read sport clubs"},
	{Name: PermSportClubCreate, Description: "Create sport clubs"},
	{Name: PermSportClubUpdate, Description: "Update sport clubs"},
	{Name: PermSportClubDelete, Description: "Delete sport clubs"},
	{Name: PermCategoryRead, Description: "Read categories"},
	{Name: PermCategoryCreate, Description: "Create categories"},
	{Name: PermCategoryUpdate, Description: "Update categories"},
	{Name: PermCategoryDelete, Description: "Delete categories"},
	{Name: PermBannerRead, Description: "Read banners (admin list)"},
	{Name: PermBannerCreate, Description: "Create banners"},
	{Name: PermBannerUpdate, Description: "Update banners"},
	{Name: PermBannerDelete, Description: "Delete banners"},
	{Name: PermSlotRead, Description: "Read slots"},
	{Name: PermSlotCreate, Description: "Create slots"},
	{Name: PermSlotUpdate, Description: "Update slots"},
	{Name: PermSlotDelete, Description: "Delete slots"},
	{Name: PermUserReadSelf, Description: "Read own profile"},
	{Name: PermUserUpdateSelf, Description: "Update own profile / avatar"},
}

var DefaultRolePermissions = map[string][]string{
	string(RoleClient): {
		PermSportClubRead,
		PermCategoryRead,
		PermSlotRead,
		PermUserReadSelf,
		PermUserUpdateSelf,
	},
	string(RolePartner): {
		PermSportClubRead, PermSportClubCreate, PermSportClubUpdate, PermSportClubDelete,
		PermCategoryRead,
		PermSlotRead, PermSlotCreate, PermSlotUpdate, PermSlotDelete,
		PermUserReadSelf, PermUserUpdateSelf,
	},
	string(RoleAdmin): {
		PermSportClubRead, PermSportClubCreate, PermSportClubUpdate, PermSportClubDelete,
		PermCategoryRead, PermCategoryCreate, PermCategoryUpdate, PermCategoryDelete,
		PermBannerRead, PermBannerCreate, PermBannerUpdate, PermBannerDelete,
		PermSlotRead, PermSlotCreate, PermSlotUpdate, PermSlotDelete,
		PermUserReadSelf, PermUserUpdateSelf,
	},
}
