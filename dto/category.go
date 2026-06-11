package dto

type CreateCategoryRequest struct {
	Name string `form:"name" validate:"required,min=2,max=100"`
}

type UpdateCategoryRequest struct {
	Name *string `form:"name" validate:"omitempty,min=2,max=100"`
}

type CategoryResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CategoryListResponse struct {
	Data  []CategoryResponse `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
