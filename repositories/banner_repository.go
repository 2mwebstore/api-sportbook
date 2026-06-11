package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type BannerRepository interface {
	Create(banner *models.Banner) error
	FindAll(page, limit int, status string) ([]models.Banner, int64, error)
	FindActive() ([]models.Banner, error)
	FindByID(id uint) (*models.Banner, error)
	Update(banner *models.Banner) error
	Delete(id uint) error
}

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) Create(banner *models.Banner) error {
	return r.db.Create(banner).Error
}

func (r *bannerRepository) FindAll(page, limit int, status string) ([]models.Banner, int64, error) {
	var banners []models.Banner
	var total int64

	q := r.db.Model(&models.Banner{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("CreatedBy").
		Order("sort_order ASC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&banners).Error
	return banners, total, err
}

func (r *bannerRepository) FindActive() ([]models.Banner, error) {
	var banners []models.Banner
	err := r.db.Where("status = ?", models.BannerStatusActive).
		Order("sort_order ASC, created_at DESC").
		Find(&banners).Error
	return banners, err
}

func (r *bannerRepository) FindByID(id uint) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.Preload("CreatedBy").First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

func (r *bannerRepository) Update(banner *models.Banner) error {
	return r.db.Save(banner).Error
}

func (r *bannerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Banner{}, id).Error
}
