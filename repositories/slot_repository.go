package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type SlotRepository interface {
	Create(slot *models.Slot) error
	FindBySportClub(sportClubID uint, page, limit int, onlyAvailable bool) ([]models.Slot, int64, error)
	FindByID(id uint) (*models.Slot, error)
	Update(slot *models.Slot) error
	Delete(id uint) error
}

type slotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) Create(slot *models.Slot) error {
	return r.db.Create(slot).Error
}

func (r *slotRepository) FindBySportClub(sportClubID uint, page, limit int, onlyAvailable bool) ([]models.Slot, int64, error) {
	var slots []models.Slot
	var total int64

	q := r.db.Model(&models.Slot{}).Where("sport_club_id = ?", sportClubID)
	if onlyAvailable {
		q = q.Where("is_available = ?", true)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("Category").Preload("CreatedBy").
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&slots).Error
	return slots, total, err
}

func (r *slotRepository) FindByID(id uint) (*models.Slot, error) {
	var slot models.Slot
	err := r.db.Preload("Category").Preload("CreatedBy").Preload("SportClub").
		First(&slot, id).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *slotRepository) Update(slot *models.Slot) error {
	return r.db.Save(slot).Error
}

func (r *slotRepository) Delete(id uint) error {
	return r.db.Delete(&models.Slot{}, id).Error
}
