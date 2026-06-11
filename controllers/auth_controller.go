package controllers

import (
	"net/http"

	"myapp/dto"
	"myapp/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService services.AuthService
	validate    *validator.Validate
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation failed",
			Details: formatValidationErrors(errs),
		})
		return
	}
	resp, err := c.authService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation failed",
			Details: formatValidationErrors(errs),
		})
		return
	}
	resp, err := c.authService.Login(&req)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() != "invalid credentials" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
