package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(cat *models.Category) error
	FindAll(page, limit int, search string) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
	FindByIDs(ids []uint) ([]models.Category, error)
	Update(cat *models.Category) error
	Delete(id uint) error
	NameExists(name string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(cat *models.Category) error {
	return r.db.Create(cat).Error
}

func (r *categoryRepository) FindAll(page, limit int, search string) ([]models.Category, int64, error) {
	var cats []models.Category
	var total int64

	q := r.db.Model(&models.Category{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Order("name ASC").Offset(offset).Limit(limit).Find(&cats).Error; err != nil {
		return nil, 0, err
	}
	return cats, total, nil
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var cat models.Category
	if err := r.db.First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindByIDs(ids []uint) ([]models.Category, error) {
	var cats []models.Category
	if len(ids) == 0 {
		return cats, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *categoryRepository) Update(cat *models.Category) error {
	return r.db.Save(cat).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) NameExists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("name = ? AND deleted_at IS NULL", name).Count(&count).Error
	return count > 0, err
}
