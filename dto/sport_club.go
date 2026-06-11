package dto

import "time"

type CreateSportClubRequest struct {
	Name        string  `form:"name"        validate:"required,min=2,max=150"`
	Lat         float64 `form:"lat"         validate:"required,latitude"`
	Lng         float64 `form:"lng"         validate:"required,longitude"`
	Location    string  `form:"location"    validate:"required,min=2,max=255"`
	IsOpen      bool    `form:"is_open"`
	OpenTime    *string `form:"open_time"   validate:"omitempty,len=5"`
	CloseTime   *string `form:"close_time"  validate:"omitempty,len=5"`
	Description *string `form:"description" validate:"omitempty,max=2000"`
	CategoryIDs []uint  `form:"category_ids"`
}

type UpdateSportClubRequest struct {
	Name          *string  `form:"name"           validate:"omitempty,min=2,max=150"`
	Lat           *float64 `form:"lat"            validate:"omitempty,latitude"`
	Lng           *float64 `form:"lng"            validate:"omitempty,longitude"`
	Location      *string  `form:"location"       validate:"omitempty,min=2,max=255"`
	IsOpen        *bool    `form:"is_open"`
	OpenTime      *string  `form:"open_time"      validate:"omitempty,len=5"`
	CloseTime     *string  `form:"close_time"     validate:"omitempty,len=5"`
	Description   *string  `form:"description"    validate:"omitempty,max=2000"`
	CategoryIDs   []uint   `form:"category_ids"`
	ReplaceImages bool     `form:"replace_images"`
}

// ── Mini embedded types ───────────────────────────────────

type CategoryMini struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"image_url,omitempty"`
}

type CreatedByMini struct {
	ID        uint    `json:"id"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Role      string  `json:"role"`
}

// ── Responses ─────────────────────────────────────────────

// SportClubResponse — list view (no slots, lighter payload)
type SportClubResponse struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Location      string         `json:"location"`
	IsOpen        bool           `json:"is_open"`
	OpenTime      *string        `json:"open_time,omitempty"`
	CloseTime     *string        `json:"close_time,omitempty"`
	Description   *string        `json:"description,omitempty"`
	ImageURLs     []string       `json:"image_urls"`
	FavoriteCount int            `json:"favorite_count"`
	Categories    []CategoryMini `json:"categories"`
	CreatedBy     *CreatedByMini `json:"created_by,omitempty"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
}

// SportClubDetailResponse — single view, includes available slots
type SportClubDetailResponse struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Location      string         `json:"location"`
	IsOpen        bool           `json:"is_open"`
	OpenTime      *string        `json:"open_time,omitempty"`
	CloseTime     *string        `json:"close_time,omitempty"`
	Description   *string        `json:"description,omitempty"`
	ImageURLs     []string       `json:"image_urls"`
	FavoriteCount int            `json:"favorite_count"`
	Categories    []CategoryMini `json:"categories"`
	Slots         []SlotMini     `json:"slots"`
	CreatedBy     *CreatedByMini `json:"created_by,omitempty"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
}

func (r *SportClubResponse) SetTimestamps(created, updated time.Time) {
	r.CreatedAt = created.Format(time.RFC3339)
	r.UpdatedAt = updated.Format(time.RFC3339)
}

type SportClubListResponse struct {
	Data  []SportClubResponse `json:"data"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}
