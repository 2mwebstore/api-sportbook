package dto

// ── Requests ──────────────────────────────────────────────

type CreateBannerRequest struct {
	Title       string  `form:"title"       validate:"required,min=2,max=150"`
	Description *string `form:"description" validate:"omitempty,max=2000"`
	LinkURL     *string `form:"link_url"    validate:"omitempty,url"`
	SortOrder   int     `form:"sort_order"`
	Status      string  `form:"status"      validate:"omitempty,oneof=active inactive"`
	// image uploaded as multipart file — key: "image"
}

type UpdateBannerRequest struct {
	Title       *string `form:"title"       validate:"omitempty,min=2,max=150"`
	Description *string `form:"description" validate:"omitempty,max=2000"`
	LinkURL     *string `form:"link_url"    validate:"omitempty,url"`
	SortOrder   *int    `form:"sort_order"`
	Status      *string `form:"status"      validate:"omitempty,oneof=active inactive"`
	// image uploaded as multipart file — key: "image" (optional, replaces existing)
}

// ── Responses ─────────────────────────────────────────────

type BannerResponse struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	Description *string        `json:"description,omitempty"`
	ImageURL    string         `json:"image_url"`
	LinkURL     *string        `json:"link_url,omitempty"`
	SortOrder   int            `json:"sort_order"`
	Status      string         `json:"status"`
	CreatedBy   *CreatedByMini `json:"created_by,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type BannerListResponse struct {
	Data  []BannerResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}
