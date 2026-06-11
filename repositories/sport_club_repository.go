package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type SportClubRepository interface {
	Create(club *models.SportClub) error
	FindAll(page, limit int, search string) ([]models.SportClub, int64, error)
	FindByID(id uint) (*models.SportClub, error)
	Update(club *models.SportClub) error
	Delete(id uint) error
	ReplaceCategories(club *models.SportClub, cats []models.Category) error
	IncrementFavorite(id uint) error
	DecrementFavorite(id uint) error
}

type sportClubRepository struct {
	db *gorm.DB
}

func NewSportClubRepository(db *gorm.DB) SportClubRepository {
	return &sportClubRepository{db: db}
}

func (r *sportClubRepository) Create(club *models.SportClub) error {
	return r.db.Create(club).Error
}

func (r *sportClubRepository) FindAll(page, limit int, search string) ([]models.SportClub, int64, error) {
	var clubs []models.SportClub
	var total int64

	q := r.db.Model(&models.SportClub{})
	if search != "" {
		q = q.Where("name LIKE ? OR location LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := q.Preload("Categories").Preload("CreatedBy").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&clubs).Error
	return clubs, total, err
}

func (r *sportClubRepository) FindByID(id uint) (*models.SportClub, error) {
	var club models.SportClub
	err := r.db.
		Preload("Categories").
		Preload("CreatedBy").
		Preload("Slots", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_available = ?", true).Order("created_at ASC")
		}).
		Preload("Slots.Category").
		Preload("Slots.CreatedBy").
		First(&club, id).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *sportClubRepository) Update(club *models.SportClub) error {
	return r.db.Save(club).Error
}

func (r *sportClubRepository) Delete(id uint) error {
	return r.db.Delete(&models.SportClub{}, id).Error
}

func (r *sportClubRepository) ReplaceCategories(club *models.SportClub, cats []models.Category) error {
	return r.db.Model(club).Association("Categories").Replace(cats)
}

func (r *sportClubRepository) IncrementFavorite(id uint) error {
	return r.db.Model(&models.SportClub{}).Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

func (r *sportClubRepository) DecrementFavorite(id uint) error {
	return r.db.Model(&models.SportClub{}).Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - 1, 0)")).Error
}
