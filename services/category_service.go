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

type CategoryService interface {
	Create(name string, image *multipart.FileHeader) (*dto.CategoryResponse, error)
	GetAll(page, limit int, search string) (*dto.CategoryListResponse, error)
	GetByID(id uint) (*dto.CategoryResponse, error)
	Update(id uint, name *string, image *multipart.FileHeader) (*dto.CategoryResponse, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(name string, image *multipart.FileHeader) (*dto.CategoryResponse, error) {
	exists, err := s.repo.NameExists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists")
	}

	cat := &models.Category{Name: name}

	if image != nil {
		url, err := utils.UploadFile(image, "categories")
		if err != nil {
			return nil, err
		}
		cat.ImageURL = &url
	}

	if err = s.repo.Create(cat); err != nil {
		return nil, err
	}
	return toCategoryResponse(cat), nil
}

func (s *categoryService) GetAll(page, limit int, search string) (*dto.CategoryListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cats, total, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	items := make([]dto.CategoryResponse, len(cats))
	for i, c := range cats {
		items[i] = *toCategoryResponse(&c)
	}
	return &dto.CategoryListResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *categoryService) GetByID(id uint) (*dto.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return toCategoryResponse(cat), nil
}

func (s *categoryService) Update(id uint, name *string, image *multipart.FileHeader) (*dto.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	if name != nil {
		if *name != cat.Name {
			exists, err := s.repo.NameExists(*name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("category name already exists")
			}
		}
		cat.Name = *name
	}

	if image != nil {
		if cat.ImageURL != nil {
			_ = utils.DeleteFile(*cat.ImageURL)
		}
		url, err := utils.UploadFile(image, "categories")
		if err != nil {
			return nil, err
		}
		cat.ImageURL = &url
	}

	if err = s.repo.Update(cat); err != nil {
		return nil, err
	}
	return toCategoryResponse(cat), nil
}

func (s *categoryService) Delete(id uint) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}
	if cat.ImageURL != nil {
		_ = utils.DeleteFile(*cat.ImageURL)
	}
	return s.repo.Delete(id)
}

func toCategoryResponse(c *models.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		ImageURL:  c.ImageURL,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}
