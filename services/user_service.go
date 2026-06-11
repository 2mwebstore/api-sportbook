package services

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"time"

	"myapp/dto"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(userID uint) (*dto.UserProfileResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserProfileResponse, error)
	UploadAvatar(userID uint, file *multipart.FileHeader) (*dto.UserProfileResponse, error)
	AddFavorite(userID, sportClubID uint) (*dto.FavoriteStatusResponse, error)
	RemoveFavorite(userID, sportClubID uint) (*dto.FavoriteStatusResponse, error)
	GetFavorites(userID uint, page, limit int) (*dto.FavoriteListResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
	scRepo   repositories.SportClubRepository
}

func NewUserService(userRepo repositories.UserRepository, scRepo repositories.SportClubRepository) UserService {
	return &userService{userRepo: userRepo, scRepo: scRepo}
}

func (s *userService) GetProfile(userID uint) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return buildProfileResponse(user), nil
}

func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	if err := s.userRepo.UpdateProfile(userID, req.FullName, req.Lat, req.Lng, req.Location); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return buildProfileResponse(user), nil
}

func (s *userService) UploadAvatar(userID uint, file *multipart.FileHeader) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user.AvatarURL != nil {
		_ = utils.DeleteFile(*user.AvatarURL)
	}
	url, err := utils.UploadFile(file, "avatars")
	if err != nil {
		return nil, err
	}
	if err = s.userRepo.UpdateAvatar(userID, url); err != nil {
		return nil, err
	}
	user.AvatarURL = &url
	return buildProfileResponse(user), nil
}

// ── Favorites ─────────────────────────────────────────────

func (s *userService) AddFavorite(userID, sportClubID uint) (*dto.FavoriteStatusResponse, error) {
	already, err := s.userRepo.IsFavorite(userID, sportClubID)
	if err != nil {
		return nil, err
	}
	if already {
		return nil, errors.New("already in favorites")
	}

	sc, err := s.scRepo.FindByID(sportClubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport club not found")
		}
		return nil, err
	}

	if err = s.userRepo.AddFavorite(userID, sportClubID); err != nil {
		return nil, err
	}
	_ = s.scRepo.IncrementFavorite(sportClubID)

	return &dto.FavoriteStatusResponse{
		SportClubID:   sportClubID,
		IsFavorited:   true,
		FavoriteCount: sc.FavoriteCount + 1,
		Message:       "added to favorites",
	}, nil
}

func (s *userService) RemoveFavorite(userID, sportClubID uint) (*dto.FavoriteStatusResponse, error) {
	exists, err := s.userRepo.IsFavorite(userID, sportClubID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("not in favorites")
	}

	sc, err := s.scRepo.FindByID(sportClubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport club not found")
		}
		return nil, err
	}

	if err = s.userRepo.RemoveFavorite(userID, sportClubID); err != nil {
		return nil, err
	}
	_ = s.scRepo.DecrementFavorite(sportClubID)

	count := sc.FavoriteCount - 1
	if count < 0 {
		count = 0
	}

	return &dto.FavoriteStatusResponse{
		SportClubID:   sportClubID,
		IsFavorited:   false,
		FavoriteCount: count,
		Message:       "removed from favorites",
	}, nil
}

func (s *userService) GetFavorites(userID uint, page, limit int) (*dto.FavoriteListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	clubs, total, err := s.userRepo.GetFavorites(userID, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.FavoriteSportClubResponse, len(clubs))
	for i, c := range clubs {
		var urls []string
		_ = json.Unmarshal([]byte(c.ImageURLs), &urls)
		if urls == nil {
			urls = []string{}
		}
		cats := make([]dto.CategoryMini, len(c.Categories))
		for j, cat := range c.Categories {
			cats[j] = dto.CategoryMini{ID: cat.ID, Name: cat.Name, ImageURL: cat.ImageURL}
		}
		items[i] = dto.FavoriteSportClubResponse{
			ID:            c.ID,
			Name:          c.Name,
			Lat:           c.Latitude,
			Lng:           c.Longitude,
			Location:      c.Location,
			IsOpen:        c.IsOpen,
			ImageURLs:     urls,
			FavoriteCount: c.FavoriteCount,
			Categories:    cats,
			CreatedAt:     c.CreatedAt.Format(time.RFC3339),
		}
	}

	return &dto.FavoriteListResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// ── helper ────────────────────────────────────────────────

func buildProfileResponse(user *models.User) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarURL:  user.AvatarURL,
		Role:       string(user.Role),
		Lat:        user.Latitude,
		Lng:        user.Longitude,
		Location:   user.Location,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}
}
