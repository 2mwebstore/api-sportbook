package services

import (
	"myapp/dto"
	"myapp/models"
	"myapp/repositories"
)

type PermissionService interface {
	GetAllPermissions() ([]dto.PermissionResponse, error)
	GetRolePermissions(role string) (*dto.RolePermissionResponse, error)
	SetRolePermissions(role string, permissions []string) (*dto.RolePermissionResponse, error)
}

type permissionService struct {
	repo repositories.PermissionRepository
}

func NewPermissionService(repo repositories.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) GetAllPermissions() ([]dto.PermissionResponse, error) {
	perms, err := s.repo.FindAllPermissions()
	if err != nil {
		return nil, err
	}
	result := make([]dto.PermissionResponse, len(perms))
	for i, p := range perms {
		result[i] = toPermResponse(p)
	}
	return result, nil
}

func (s *permissionService) GetRolePermissions(role string) (*dto.RolePermissionResponse, error) {
	rps, err := s.repo.GetRolePermissions(role)
	if err != nil {
		return nil, err
	}
	perms := make([]dto.PermissionResponse, len(rps))
	for i, rp := range rps {
		perms[i] = toPermResponse(rp.Permission)
	}
	return &dto.RolePermissionResponse{Role: role, Permissions: perms}, nil
}

func (s *permissionService) SetRolePermissions(role string, permissions []string) (*dto.RolePermissionResponse, error) {
	if err := s.repo.SetPermissionsForRole(role, permissions); err != nil {
		return nil, err
	}
	return s.GetRolePermissions(role)
}

func toPermResponse(p models.Permission) dto.PermissionResponse {
	return dto.PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}
