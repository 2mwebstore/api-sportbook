package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	// Role permissions
	GetPermissionsForRole(role string) ([]string, error)
	SetPermissionsForRole(role string, permissionNames []string) error
	HasPermission(role, permissionName string) (bool, error)

	// Permission CRUD
	FindAllPermissions() ([]models.Permission, error)
	FindPermissionByName(name string) (*models.Permission, error)

	// Role-permission list (for admin UI)
	GetRolePermissions(role string) ([]models.RolePermission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetPermissionsForRole(role string) ([]string, error) {
	var rps []models.RolePermission
	err := r.db.Preload("Permission").Where("role = ?", role).Find(&rps).Error
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(rps))
	for _, rp := range rps {
		names = append(names, rp.Permission.Name)
	}
	return names, nil
}

func (r *permissionRepository) HasPermission(role, permissionName string) (bool, error) {
	var count int64
	err := r.db.Table("role_permissions rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("rp.role = ? AND p.name = ? AND p.deleted_at IS NULL", role, permissionName).
		Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) SetPermissionsForRole(role string, permissionNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove all existing permissions for this role
		if err := tx.Where("role = ?", role).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		// Re-insert
		for _, name := range permissionNames {
			var perm models.Permission
			if err := tx.Where("name = ?", name).First(&perm).Error; err != nil {
				return err
			}
			rp := models.RolePermission{Role: role, PermissionID: perm.ID}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *permissionRepository) FindAllPermissions() ([]models.Permission, error) {
	var perms []models.Permission
	err := r.db.Order("name ASC").Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindPermissionByName(name string) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.Where("name = ?", name).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) GetRolePermissions(role string) ([]models.RolePermission, error) {
	var rps []models.RolePermission
	err := r.db.Preload("Permission").Where("role = ?", role).Find(&rps).Error
	return rps, err
}
