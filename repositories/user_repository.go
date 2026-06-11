package repositories

import (
	"myapp/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	UpdateLastLogin(id uint) error
	UpdateAvatar(id uint, url string) error
	UpdateLocation(id uint, lat, lng *float64, location *string) error
	UpdateProfile(id uint, fullName *string, lat, lng *float64, location *string) error
	EmailExists(email string) (bool, error)
	PhoneExists(phone string) (bool, error)
	// Favorites
	AddFavorite(userID, sportClubID uint) error
	RemoveFavorite(userID, sportClubID uint) error
	IsFavorite(userID, sportClubID uint) (bool, error)
	GetFavorites(userID uint, page, limit int) ([]models.SportClub, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

func (r *userRepository) UpdateAvatar(id uint, url string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("avatar_url", url).Error
}

func (r *userRepository) UpdateLocation(id uint, lat, lng *float64, location *string) error {
	updates := map[string]interface{}{}
	if lat != nil {
		updates["latitude"] = lat
	}
	if lng != nil {
		updates["longitude"] = lng
	}
	if location != nil {
		updates["location"] = location
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) UpdateProfile(id uint, fullName *string, lat, lng *float64, location *string) error {
	updates := map[string]interface{}{}
	if fullName != nil {
		updates["full_name"] = fullName
	}
	if lat != nil {
		updates["latitude"] = lat
	}
	if lng != nil {
		updates["longitude"] = lng
	}
	if location != nil {
		updates["location"] = location
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) PhoneExists(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("phone = ? AND deleted_at IS NULL", phone).Count(&count).Error
	return count > 0, err
}

// ── Favorites ─────────────────────────────────────────────

func (r *userRepository) AddFavorite(userID, sportClubID uint) error {
	user := models.User{}
	user.ID = userID
	sc := models.SportClub{}
	sc.ID = sportClubID
	return r.db.Model(&user).Association("Favorites").Append(&sc)
}

func (r *userRepository) RemoveFavorite(userID, sportClubID uint) error {
	user := models.User{}
	user.ID = userID
	sc := models.SportClub{}
	sc.ID = sportClubID
	return r.db.Model(&user).Association("Favorites").Delete(&sc)
}

func (r *userRepository) IsFavorite(userID, sportClubID uint) (bool, error) {
	var count int64
	err := r.db.Table("user_favorites").
		Where("user_id = ? AND sport_club_id = ?", userID, sportClubID).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) GetFavorites(userID uint, page, limit int) ([]models.SportClub, int64, error) {
	var user models.User
	user.ID = userID

	total := int64(r.db.Model(&user).Association("Favorites").Count())

	offset := (page - 1) * limit
	var clubs []models.SportClub
	err := r.db.Model(&user).
		Preload("Categories").
		Association("Favorites").
		Find(&clubs)
	if err != nil {
		return nil, 0, err
	}

	// Manual pagination on the slice
	start := offset
	if start > len(clubs) {
		start = len(clubs)
	}
	end := start + limit
	if end > len(clubs) {
		end = len(clubs)
	}

	return clubs[start:end], total, nil
}
