package services

import (
	"errors"

	"myapp/dto"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils"
	"myapp/validators"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if err := validators.ValidateRegister(req); err != nil {
		return nil, err
	}

	if req.Email != nil {
		exists, err := s.userRepo.EmailExists(*req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email is already registered")
		}
	}
	if req.Phone != nil {
		exists, err := s.userRepo.PhoneExists(*req.Phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("phone number is already registered")
		}
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := models.RoleClient
	if req.Role == string(models.RolePartner) {
		role = models.RolePartner
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashed,
		Role:     role,
	}
	if err = s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if err := validators.ValidateLogin(req); err != nil {
		return nil, err
	}

	var user *models.User
	var err error

	switch {
	case req.Email != nil:
		user, err = s.userRepo.FindByEmail(*req.Email)
	case req.Phone != nil:
		user, err = s.userRepo.FindByPhone(*req.Phone)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return s.buildAuthResponse(user)
}

func (s *authService) buildAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	phone := ""
	if user.Phone != nil {
		phone = *user.Phone
	}

	access, refresh, expiresIn, err := utils.GenerateTokenPair(user.ID, string(user.Role), email, phone)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User: dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Phone:      user.Phone,
			AvatarURL:  user.AvatarURL,
			Role:       string(user.Role),
			IsVerified: user.IsVerified,
			IsActive:   user.IsActive,
		},
	}, nil
}
