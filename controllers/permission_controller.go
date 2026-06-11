package controllers

import (
	"net/http"

	"myapp/dto"
	"myapp/models"
	"myapp/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PermissionController struct {
	service  services.PermissionService
	validate *validator.Validate
}

func NewPermissionController(service services.PermissionService) *PermissionController {
	return &PermissionController{service: service, validate: validator.New()}
}

// GetAllPermissions — GET /api/v1/admin/permissions
// Lists every permission that exists in the system.
func (c *PermissionController) GetAllPermissions(ctx *gin.Context) {
	perms, err := c.service.GetAllPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": perms})
}

// GetRolePermissions — GET /api/v1/admin/permissions/roles/:role
// Returns permissions currently assigned to a role.
func (c *PermissionController) GetRolePermissions(ctx *gin.Context) {
	role := ctx.Param("role")
	if !validRole(role) {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid role: must be client, partner or admin"})
		return
	}
	resp, err := c.service.GetRolePermissions(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// SetRolePermissions — PUT /api/v1/admin/permissions/roles/:role
// Replaces all permissions for a role. Send the full desired list.
func (c *PermissionController) SetRolePermissions(ctx *gin.Context) {
	role := ctx.Param("role")
	if !validRole(role) {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid role: must be client, partner or admin"})
		return
	}

	var req dto.SetRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation failed", Details: formatValidationErrors(errs),
		})
		return
	}

	resp, err := c.service.SetRolePermissions(role, req.Permissions)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func validRole(role string) bool {
	return role == string(models.RoleClient) ||
		role == string(models.RolePartner) ||
		role == string(models.RoleAdmin)
}
