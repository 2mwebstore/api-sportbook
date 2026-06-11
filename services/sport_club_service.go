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

type SportClubService interface {
	Create(userID uint, req *dto.CreateSportClubRequest, images []*multipart.FileHeader) (*dto.SportClubResponse, error)
	GetAll(page, limit int, search string) (*dto.SportClubListResponse, error)
	GetByID(id uint) (*dto.SportClubDetailResponse, error)
	Update(id uint, req *dto.UpdateSportClubRequest, newImages []*multipart.FileHeader) (*dto.SportClubResponse, error)
	Delete(id uint) error
}

type sportClubService struct {
	repo    repositories.SportClubRepository
	catRepo repositories.CategoryRepository
}

func NewSportClubService(repo repositories.SportClubRepository, catRepo repositories.CategoryRepository) SportClubService {
	return &sportClubService{repo: repo, catRepo: catRepo}
}

func (s *sportClubService) Create(userID uint, req *dto.CreateSportClubRequest, images []*multipart.FileHeader) (*dto.SportClubResponse, error) {
	imageURLs, err := utils.UploadMultipleFiles(images, "sport-clubs")
	if err != nil {
		return nil, err
	}
	imageJSON, _ := marshalURLs(imageURLs)

	club := &models.SportClub{
		Name:        req.Name,
		Latitude:    req.Lat,
		Longitude:   req.Lng,
		Location:    req.Location,
		IsOpen:      req.IsOpen,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		Description: req.Description,
		ImageURLs:   imageJSON,
		CreatedByID: &userID,
	}

	if len(req.CategoryIDs) > 0 {
		cats, err := s.catRepo.FindByIDs(req.CategoryIDs)
		if err != nil {
			return nil, err
		}
		club.Categories = cats
	}

	if err = s.repo.Create(club); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(club.ID)
	if err != nil {
		return nil, err
	}
	return toSportClubResponse(created), nil
}

func (s *sportClubService) GetAll(page, limit int, search string) (*dto.SportClubListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	clubs, total, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SportClubResponse, len(clubs))
	for i, c := range clubs {
		items[i] = *toSportClubResponse(&c)
	}
	return &dto.SportClubListResponse{Data: items, Total: total, Page: page, Limit: limit}, nil
}

// GetByID returns the full detail including available slots.
func (s *sportClubService) GetByID(id uint) (*dto.SportClubDetailResponse, error) {
	club, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport club not found")
		}
		return nil, err
	}
	return toSportClubDetailResponse(club), nil
}

func (s *sportClubService) Update(id uint, req *dto.UpdateSportClubRequest, newImages []*multipart.FileHeader) (*dto.SportClubResponse, error) {
	club, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sport club not found")
		}
		return nil, err
	}

	if req.Name != nil {
		club.Name = *req.Name
	}
	if req.Lat != nil {
		club.Latitude = *req.Lat
	}
	if req.Lng != nil {
		club.Longitude = *req.Lng
	}
	if req.Location != nil {
		club.Location = *req.Location
	}
	if req.IsOpen != nil {
		club.IsOpen = *req.IsOpen
	}
	if req.OpenTime != nil {
		club.OpenTime = req.OpenTime
	}
	if req.CloseTime != nil {
		club.CloseTime = req.CloseTime
	}
	if req.Description != nil {
		club.Description = req.Description
	}

	if len(newImages) > 0 {
		if req.ReplaceImages {
			for _, u := range unmarshalURLs(club.ImageURLs) {
				_ = utils.DeleteFile(u)
			}
			uploaded, err := utils.UploadMultipleFiles(newImages, "sport-clubs")
			if err != nil {
				return nil, err
			}
			club.ImageURLs, _ = marshalURLs(uploaded)
		} else {
			existing := unmarshalURLs(club.ImageURLs)
			uploaded, err := utils.UploadMultipleFiles(newImages, "sport-clubs")
			if err != nil {
				return nil, err
			}
			club.ImageURLs, _ = marshalURLs(append(existing, uploaded...))
		}
	}

	if len(req.CategoryIDs) > 0 {
		cats, err := s.catRepo.FindByIDs(req.CategoryIDs)
		if err != nil {
			return nil, err
		}
		if err = s.repo.ReplaceCategories(club, cats); err != nil {
			return nil, err
		}
		club.Categories = cats
	}

	if err = s.repo.Update(club); err != nil {
		return nil, err
	}
	return toSportClubResponse(club), nil
}

func (s *sportClubService) Delete(id uint) error {
	club, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("sport club not found")
		}
		return err
	}
	for _, u := range unmarshalURLs(club.ImageURLs) {
		_ = utils.DeleteFile(u)
	}
	return s.repo.Delete(id)
}

// ── helpers ───────────────────────────────────────────────

func marshalURLs(urls []string) (string, error) {
	if urls == nil {
		urls = []string{}
	}
	b, err := json.Marshal(urls)
	return string(b), err
}

func unmarshalURLs(raw string) []string {
	var urls []string
	if err := json.Unmarshal([]byte(raw), &urls); err != nil {
		return []string{}
	}
	return urls
}

func toSportClubResponse(c *models.SportClub) *dto.SportClubResponse {
	cats := make([]dto.CategoryMini, len(c.Categories))
	for i, cat := range c.Categories {
		cats[i] = dto.CategoryMini{ID: cat.ID, Name: cat.Name, ImageURL: cat.ImageURL}
	}
	resp := &dto.SportClubResponse{
		ID:            c.ID,
		Name:          c.Name,
		Lat:           c.Latitude,
		Lng:           c.Longitude,
		Location:      c.Location,
		IsOpen:        c.IsOpen,
		OpenTime:      c.OpenTime,
		CloseTime:     c.CloseTime,
		Description:   c.Description,
		ImageURLs:     unmarshalURLs(c.ImageURLs),
		FavoriteCount: c.FavoriteCount,
		Categories:    cats,
		CreatedAt:     c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     c.UpdatedAt.Format(time.RFC3339),
	}
	if c.CreatedBy != nil {
		resp.CreatedBy = &dto.CreatedByMini{
			ID: c.CreatedBy.ID, FullName: c.CreatedBy.FullName,
			AvatarURL: c.CreatedBy.AvatarURL, Role: string(c.CreatedBy.Role),
		}
	}
	return resp
}

func toSportClubDetailResponse(c *models.SportClub) *dto.SportClubDetailResponse {
	cats := make([]dto.CategoryMini, len(c.Categories))
	for i, cat := range c.Categories {
		cats[i] = dto.CategoryMini{ID: cat.ID, Name: cat.Name, ImageURL: cat.ImageURL}
	}

	slots := make([]dto.SlotMini, len(c.Slots))
	for i, sl := range c.Slots {
		slots[i] = dto.SlotMini{
			ID:          sl.ID,
			Name:        sl.Name,
			ImageURL:    sl.ImageURL,
			Price:       sl.Price,
			Capacity:    sl.Capacity,
			IsAvailable: sl.IsAvailable,
		}
	}

	resp := &dto.SportClubDetailResponse{
		ID:            c.ID,
		Name:          c.Name,
		Lat:           c.Latitude,
		Lng:           c.Longitude,
		Location:      c.Location,
		IsOpen:        c.IsOpen,
		OpenTime:      c.OpenTime,
		CloseTime:     c.CloseTime,
		Description:   c.Description,
		ImageURLs:     unmarshalURLs(c.ImageURLs),
		FavoriteCount: c.FavoriteCount,
		Categories:    cats,
		Slots:         slots,
		CreatedAt:     c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     c.UpdatedAt.Format(time.RFC3339),
	}
	if c.CreatedBy != nil {
		resp.CreatedBy = &dto.CreatedByMini{
			ID: c.CreatedBy.ID, FullName: c.CreatedBy.FullName,
			AvatarURL: c.CreatedBy.AvatarURL, Role: string(c.CreatedBy.Role),
		}
	}
	return resp
}
