package dto

// ── Update profile ────────────────────────────────────────

type UpdateProfileRequest struct {
	FullName *string  `json:"full_name" validate:"omitempty,min=2,max=100"`
	Lat      *float64 `json:"lat"       validate:"omitempty,latitude"`
	Lng      *float64 `json:"lng"       validate:"omitempty,longitude"`
	Location *string  `json:"location"  validate:"omitempty,max=255"`
}

// ── Favorite ──────────────────────────────────────────────

type FavoriteSportClubResponse struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Location      string         `json:"location"`
	IsOpen        bool           `json:"is_open"`
	ImageURLs     []string       `json:"image_urls"`
	FavoriteCount int            `json:"favorite_count"`
	Categories    []CategoryMini `json:"categories"`
	CreatedAt     string         `json:"created_at"`
}

type FavoriteListResponse struct {
	Data  []FavoriteSportClubResponse `json:"data"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Limit int                         `json:"limit"`
}

type FavoriteStatusResponse struct {
	SportClubID   uint   `json:"sport_club_id"`
	IsFavorited   bool   `json:"is_favorited"`
	FavoriteCount int    `json:"favorite_count"`
	Message       string `json:"message"`
}
