package dto

// ── Requests ──────────────────────────────────────────────

type CreateSlotRequest struct {
	SportClubID uint    `form:"sport_club_id" validate:"required"`
	Name        string  `form:"name"          validate:"required,min=2,max=150"`
	Description *string `form:"description"   validate:"omitempty,max=2000"`
	Price       float64 `form:"price"         validate:"min=0"`
	Capacity    int     `form:"capacity"      validate:"min=1"`
	IsAvailable bool    `form:"is_available"`
	CategoryID  *uint   `form:"category_id"`
	// image uploaded as multipart file — key: "image"
}

type UpdateSlotRequest struct {
	SportClubID *uint    `form:"sport_club_id" validate:"omitempty"`
	Name        *string  `form:"name"          validate:"omitempty,min=2,max=150"`
	Description *string  `form:"description"   validate:"omitempty,max=2000"`
	Price       *float64 `form:"price"         validate:"omitempty,min=0"`
	Capacity    *int     `form:"capacity"      validate:"omitempty,min=1"`
	IsAvailable *bool    `form:"is_available"`
	CategoryID  *uint    `form:"category_id"`
	// image uploaded as multipart file — key: "image" (optional, replaces existing)
}

// ── Responses ─────────────────────────────────────────────

type SlotMini struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	ImageURL    *string `json:"image_url,omitempty"`
	Price       float64 `json:"price"`
	Capacity    int     `json:"capacity"`
	IsAvailable bool    `json:"is_available"`
}

type SlotResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	ImageURL    *string        `json:"image_url,omitempty"`
	Description *string        `json:"description,omitempty"`
	Price       float64        `json:"price"`
	Capacity    int            `json:"capacity"`
	IsAvailable bool           `json:"is_available"`
	SportClubID uint           `json:"sport_club_id"`
	Category    *CategoryMini  `json:"category,omitempty"`
	CreatedBy   *CreatedByMini `json:"created_by,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type SlotListResponse struct {
	Data  []SlotResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
