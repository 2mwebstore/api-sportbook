package validators

import (
	"errors"

	"myapp/dto"
)

// ValidateRegister enforces business rules on top of struct-tag validation.
func ValidateRegister(req *dto.RegisterRequest) error {
	if req.Email == nil && req.Phone == nil {
		return errors.New("at least one of email or phone is required")
	}
	return nil
}

// ValidateLogin ensures exactly one identifier is supplied.
func ValidateLogin(req *dto.LoginRequest) error {
	if req.Email == nil && req.Phone == nil {
		return errors.New("email or phone is required")
	}
	if req.Email != nil && req.Phone != nil {
		return errors.New("provide either email or phone, not both")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
