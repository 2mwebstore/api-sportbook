package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"myapp/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	logLevel := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// ── Connection pooling (important for Vercel serverless) ──
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(5) // GoDaddy shared hosting has low connection limits
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// ──────────────────────────────────────────────────────────

	DB = db
	log.Println("database connection established")
	return db
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Category{},
		&models.SportClub{},
		&models.Slot{},
		&models.Banner{},
	); err != nil {
		log.Fatalf("auto-migration failed: %v", err)
	}
	log.Println("database migration completed")
}
