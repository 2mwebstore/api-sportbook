package dto

// ── Refresh ───────────────────────────────────────────────

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ── Register ──────────────────────────────────────────────

type RegisterRequest struct {
	FullName string  `json:"full_name" validate:"required,min=2,max=100"`
	Email    *string `json:"email"     validate:"omitempty,email"`
	Phone    *string `json:"phone"     validate:"omitempty,e164"`
	Password string  `json:"password"  validate:"required,min=8"`
	Role     string  `json:"role"      validate:"omitempty,oneof=client partner"`
}

// ── Login ─────────────────────────────────────────────────

type LoginRequest struct {
	Email    *string `json:"email"    validate:"omitempty,email"`
	Phone    *string `json:"phone"    validate:"omitempty,e164"`
	Password string  `json:"password" validate:"required"`
}

// ── Responses ─────────────────────────────────────────────

type UserResponse struct {
	ID         uint    `json:"id"`
	FullName   string  `json:"full_name"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	Role       string  `json:"role"`
	IsVerified bool    `json:"is_verified"`
	IsActive   bool    `json:"is_active"`
}

type UserProfileResponse struct {
	ID         uint     `json:"id"`
	FullName   string   `json:"full_name"`
	Email      *string  `json:"email,omitempty"`
	Phone      *string  `json:"phone,omitempty"`
	AvatarURL  *string  `json:"avatar_url,omitempty"`
	Role       string   `json:"role"`
	Lat        *float64 `json:"lat,omitempty"`
	Lng        *float64 `json:"lng,omitempty"`
	Location   *string  `json:"location,omitempty"`
	IsVerified bool     `json:"is_verified"`
	IsActive   bool     `json:"is_active"`
	CreatedAt  string   `json:"created_at"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
