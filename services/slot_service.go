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

type SlotService interface {
	Create(userID uint, req *dto.CreateSlotRequest, image *multipart.FileHeader) (*dto.SlotResponse, error)
	GetBySportClub(sportClubID uint, page, limit int, onlyAvailable bool) (*dto.SlotListResponse, error)
	GetByID(id uint) (*dto.SlotResponse, error)
	Update(id uint, req *dto.UpdateSlotRequest, image *multipart.FileHeader) (*dto.SlotResponse, error)
	Delete(id uint) error
}

type slotService struct {
	repo    repositories.SlotRepository
	scRepo  repositories.SportClubRepository
	catRepo repositories.CategoryRepository
}

func NewSlotService(
	repo repositories.SlotRepository,
	scRepo repositories.SportClubRepository,
	catRepo repositories.CategoryRepository,
) SlotService {
	return &slotService{repo: repo, scRepo: scRepo, catRepo: catRepo}
}

func (s *slotService) Create(userID uint, req *dto.CreateSlotRequest, image *multipart.FileHeader) (*dto.SlotResponse, error) {
	// Verify sport club exists
	if _, err := s.scRepo.FindByID(req.SportClubID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport club not found")
		}
		return nil, err
	}

	slot := &models.Slot{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Capacity:    req.Capacity,
		IsAvailable: req.IsAvailable,
		SportClubID: req.SportClubID,
		CategoryID:  req.CategoryID,
		CreatedByID: &userID,
	}

	if image != nil {
		url, err := utils.UploadFile(image, "slots")
		if err != nil {
			return nil, err
		}
		slot.ImageURL = &url
	}

	if err := s.repo.Create(slot); err != nil {
		return nil, err
	}

	// Reload to get relations
	created, err := s.repo.FindByID(slot.ID)
	if err != nil {
		return nil, err
	}
	return toSlotResponse(created), nil
}

func (s *slotService) GetBySportClub(sportClubID uint, page, limit int, onlyAvailable bool) (*dto.SlotListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	slots, total, err := s.repo.FindBySportClub(sportClubID, page, limit, onlyAvailable)
	if err != nil {
		return nil, err
	}

	items := make([]dto.SlotResponse, len(slots))
	for i, sl := range slots {
		items[i] = *toSlotResponse(&sl)
	}
	return &dto.SlotListResponse{Data: items, Total: total, Page: page, Limit: limit}, nil
}

func (s *slotService) GetByID(id uint) (*dto.SlotResponse, error) {
	slot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("slot not found")
		}
		return nil, err
	}
	return toSlotResponse(slot), nil
}

func (s *slotService) Update(id uint, req *dto.UpdateSlotRequest, image *multipart.FileHeader) (*dto.SlotResponse, error) {
	slot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("slot not found")
		}
		return nil, err
	}

	if req.SportClubID != nil {
		// Verify new sport club exists before reassigning
		if _, err := s.scRepo.FindByID(*req.SportClubID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("sport club not found")
			}
			return nil, err
		}
		slot.SportClubID = *req.SportClubID
	}
	if req.Name != nil {
		slot.Name = *req.Name
	}
	if req.Description != nil {
		slot.Description = req.Description
	}
	if req.Price != nil {
		slot.Price = *req.Price
	}
	if req.Capacity != nil {
		slot.Capacity = *req.Capacity
	}
	if req.IsAvailable != nil {
		slot.IsAvailable = *req.IsAvailable
	}
	if req.CategoryID != nil {
		slot.CategoryID = req.CategoryID
	}

	if image != nil {
		if slot.ImageURL != nil {
			_ = utils.DeleteFile(*slot.ImageURL)
		}
		url, err := utils.UploadFile(image, "slots")
		if err != nil {
			return nil, err
		}
		slot.ImageURL = &url
	}

	if err = s.repo.Update(slot); err != nil {
		return nil, err
	}

	// Reload to reflect latest relations
	updated, err := s.repo.FindByID(slot.ID)
	if err != nil {
		return nil, err
	}
	return toSlotResponse(updated), nil
}

func (s *slotService) Delete(id uint) error {
	slot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("slot not found")
		}
		return err
	}
	if slot.ImageURL != nil {
		_ = utils.DeleteFile(*slot.ImageURL)
	}
	return s.repo.Delete(id)
}

// ── helper ────────────────────────────────────────────────

func toSlotResponse(s *models.Slot) *dto.SlotResponse {
	resp := &dto.SlotResponse{
		ID:          s.ID,
		Name:        s.Name,
		ImageURL:    s.ImageURL,
		Description: s.Description,
		Price:       s.Price,
		Capacity:    s.Capacity,
		IsAvailable: s.IsAvailable,
		SportClubID: s.SportClubID,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
	if s.Category != nil {
		resp.Category = &dto.CategoryMini{
			ID:       s.Category.ID,
			Name:     s.Category.Name,
			ImageURL: s.Category.ImageURL,
		}
	}
	if s.CreatedBy != nil {
		resp.CreatedBy = &dto.CreatedByMini{
			ID:        s.CreatedBy.ID,
			FullName:  s.CreatedBy.FullName,
			AvatarURL: s.CreatedBy.AvatarURL,
			Role:      string(s.CreatedBy.Role),
		}
	}
	return resp
}
