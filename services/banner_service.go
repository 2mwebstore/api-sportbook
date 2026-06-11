package services

import (
	"errors"
	"mime/multipart"
	"time"

	"myapp/dto"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils"

	"gorm.io/gorm"
)

type BannerService interface {
	Create(userID uint, req *dto.CreateBannerRequest, image *multipart.FileHeader) (*dto.BannerResponse, error)
	GetAll(page, limit int, status string) (*dto.BannerListResponse, error)
	GetActive() ([]dto.BannerResponse, error)
	GetByID(id uint) (*dto.BannerResponse, error)
	Update(id uint, req *dto.UpdateBannerRequest, image *multipart.FileHeader) (*dto.BannerResponse, error)
	Delete(id uint) error
}

type bannerService struct {
	repo repositories.BannerRepository
}

func NewBannerService(repo repositories.BannerRepository) BannerService {
	return &bannerService{repo: repo}
}

func (s *bannerService) Create(userID uint, req *dto.CreateBannerRequest, image *multipart.FileHeader) (*dto.BannerResponse, error) {
	if image == nil {
		return nil, errors.New("banner image is required")
	}

	imageURL, err := utils.UploadFile(image, "banners")
	if err != nil {
		return nil, err
	}

	status := models.BannerStatusActive
	if req.Status == string(models.BannerStatusInactive) {
		status = models.BannerStatusInactive
	}

	banner := &models.Banner{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    imageURL,
		LinkURL:     req.LinkURL,
		SortOrder:   req.SortOrder,
		Status:      status,
		CreatedByID: &userID,
	}

	if err = s.repo.Create(banner); err != nil {
		return nil, err
	}

	// reload to get CreatedBy populated
	created, err := s.repo.FindByID(banner.ID)
	if err != nil {
		return nil, err
	}
	return toBannerResponse(created), nil
}

func (s *bannerService) GetAll(page, limit int, status string) (*dto.BannerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	banners, total, err := s.repo.FindAll(page, limit, status)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BannerResponse, len(banners))
	for i, b := range banners {
		items[i] = *toBannerResponse(&b)
	}
	return &dto.BannerListResponse{Data: items, Total: total, Page: page, Limit: limit}, nil
}

func (s *bannerService) GetActive() ([]dto.BannerResponse, error) {
	banners, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}
	items := make([]dto.BannerResponse, len(banners))
	for i, b := range banners {
		items[i] = *toBannerResponse(&b)
	}
	return items, nil
}

func (s *bannerService) GetByID(id uint) (*dto.BannerResponse, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("banner not found")
		}
		return nil, err
	}
	return toBannerResponse(banner), nil
}

func (s *bannerService) Update(id uint, req *dto.UpdateBannerRequest, image *multipart.FileHeader) (*dto.BannerResponse, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("banner not found")
		}
		return nil, err
	}

	if req.Title != nil {
		banner.Title = *req.Title
	}
	if req.Description != nil {
		banner.Description = req.Description
	}
	if req.LinkURL != nil {
		banner.LinkURL = req.LinkURL
	}
	if req.SortOrder != nil {
		banner.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		banner.Status = models.BannerStatus(*req.Status)
	}

	// Replace image if new one provided
	if image != nil {
		_ = utils.DeleteFile(banner.ImageURL)
		newURL, err := utils.UploadFile(image, "banners")
		if err != nil {
			return nil, err
		}
		banner.ImageURL = newURL
	}

	if err = s.repo.Update(banner); err != nil {
		return nil, err
	}
	return toBannerResponse(banner), nil
}

func (s *bannerService) Delete(id uint) error {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("banner not found")
		}
		return err
	}
	_ = utils.DeleteFile(banner.ImageURL)
	return s.repo.Delete(id)
}

// ── helper ────────────────────────────────────────────────

func toBannerResponse(b *models.Banner) *dto.BannerResponse {
	resp := &dto.BannerResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		ImageURL:    b.ImageURL,
		LinkURL:     b.LinkURL,
		SortOrder:   b.SortOrder,
		Status:      string(b.Status),
		CreatedAt:   b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
	}
	if b.CreatedBy != nil {
		resp.CreatedBy = &dto.CreatedByMini{
			ID:        b.CreatedBy.ID,
			FullName:  b.CreatedBy.FullName,
			AvatarURL: b.CreatedBy.AvatarURL,
			Role:      string(b.CreatedBy.Role),
		}
	}
	return resp
}
