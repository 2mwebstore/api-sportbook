package database

import (
	"log"
	"os"

	"myapp/models"
	"myapp/utils"

	"gorm.io/gorm"
)

// Seed runs all seeders. Safe to call on every startup — all ops are idempotent.
func Seed(db *gorm.DB) {
	seedPermissions(db)
	seedRolePermissions(db)
	seedOwnerUser(db)
}

// seedPermissions inserts any missing permissions from the master list.
func seedPermissions(db *gorm.DB) {
	for _, p := range models.AllPermissions {
		var existing models.Permission
		result := db.Where("name = ?", p.Name).First(&existing)
		if result.Error != nil {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("[seeder] failed to create permission %s: %v", p.Name, err)
			} else {
				log.Printf("[seeder] permission created: %s", p.Name)
			}
		}
	}
}

// seedRolePermissions assigns default permissions to each role if not already set.
func seedRolePermissions(db *gorm.DB) {
	for role, permNames := range models.DefaultRolePermissions {
		var count int64
		db.Model(&models.RolePermission{}).Where("role = ?", role).Count(&count)
		if count > 0 {
			continue // already seeded for this role
		}
		for _, name := range permNames {
			var perm models.Permission
			if err := db.Where("name = ?", name).First(&perm).Error; err != nil {
				log.Printf("[seeder] permission not found: %s", name)
				continue
			}
			rp := models.RolePermission{Role: role, PermissionID: perm.ID}
			if err := db.Create(&rp).Error; err != nil {
				log.Printf("[seeder] failed to assign %s → %s: %v", role, name, err)
			}
		}
		log.Printf("[seeder] role permissions seeded for: %s", role)
	}
}

// seedOwnerUser creates user ID=1 (super admin/owner) from env vars if not present.
func seedOwnerUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("id = 1").Count(&count)
	if count > 0 {
		return // owner already exists
	}

	email := os.Getenv("OWNER_EMAIL")
	password := os.Getenv("OWNER_PASSWORD")
	fullName := os.Getenv("OWNER_NAME")

	if email == "" {
		email = "owner@myapp.com"
	}
	if password == "" {
		password = "Owner@123456"
	}
	if fullName == "" {
		fullName = "System Owner"
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("[seeder] failed to hash owner password: %v", err)
	}

	owner := models.User{
		FullName:   fullName,
		Email:      &email,
		Password:   hashed,
		Role:       models.RoleAdmin,
		IsVerified: true,
		IsActive:   true,
	}

	// Force ID = 1
	if err := db.Create(&owner).Error; err != nil {
		log.Fatalf("[seeder] failed to create owner user: %v", err)
	}
	log.Printf("[seeder] owner user created — email: %s  password: %s", email, password)
}
