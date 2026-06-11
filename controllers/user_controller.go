package controllers

import (
	"net/http"
	"strconv"

	"myapp/dto"
	"myapp/middleware"
	"myapp/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService, validate: validator.New()}
}

// GetProfile — GET /api/v1/users/me
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)
	profile, err := c.userService.GetProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

// UpdateProfile — PUT /api/v1/users/me
// JSON body: full_name, lat, lng, location (all optional)
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)

	var req dto.UpdateProfileRequest
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

	profile, err := c.userService.UpdateProfile(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

// UploadAvatar — POST /api/v1/users/me/avatar
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "avatar file is required"})
		return
	}
	ct := file.Header.Get("Content-Type")
	allowed := map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true, "image/gif": true}
	if !allowed[ct] {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "unsupported image type; use jpeg, png, webp or gif"})
		return
	}
	if file.Size > 5<<20 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "avatar must be under 5 MB"})
		return
	}

	profile, err := c.userService.UploadAvatar(userID, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

// ── Favorites ─────────────────────────────────────────────

// GetFavorites — GET /api/v1/users/me/favorites?page=1&limit=10
func (c *UserController) GetFavorites(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.userService.GetFavorites(userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// AddFavorite — POST /api/v1/users/me/favorites/:sport_club_id
func (c *UserController) AddFavorite(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)
	sportClubID, err := parseSportClubID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid sport_club_id"})
		return
	}

	resp, err := c.userService.AddFavorite(userID, sportClubID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "sport club not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// RemoveFavorite — DELETE /api/v1/users/me/favorites/:sport_club_id
func (c *UserController) RemoveFavorite(ctx *gin.Context) {
	userID := ctx.GetUint(middleware.ContextUserID)
	sportClubID, err := parseSportClubID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid sport_club_id"})
		return
	}

	resp, err := c.userService.RemoveFavorite(userID, sportClubID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "not in favorites" || err.Error() == "sport club not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func parseSportClubID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("sport_club_id"), 10, 64)
	return uint(id), err
}
