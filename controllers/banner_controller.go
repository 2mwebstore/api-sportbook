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

type BannerController struct {
	service  services.BannerService
	validate *validator.Validate
}

func NewBannerController(service services.BannerService) *BannerController {
	return &BannerController{service: service, validate: validator.New()}
}

// GetActive — GET /api/v1/banners/active
// Public: returns only active banners ordered by sort_order.
func (c *BannerController) GetActive(ctx *gin.Context) {
	banners, err := c.service.GetActive()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": banners})
}

// GetAll — GET /api/v1/admin/banners?page=1&limit=10&status=active
// Admin: paginated list, optional status filter.
func (c *BannerController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	status := ctx.Query("status") // "active" | "inactive" | "" (all)

	resp, err := c.service.GetAll(page, limit, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetByID — GET /api/v1/admin/banners/:id
func (c *BannerController) GetByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	resp, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Create — POST /api/v1/admin/banners  (multipart/form-data)
// Fields : title, description, link_url, sort_order, status
// File   : image (required)
func (c *BannerController) Create(ctx *gin.Context) {
	var req dto.CreateBannerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation failed", Details: formatValidationErrors(errs),
		})
		return
	}

	image, _ := ctx.FormFile("image")
	userID := ctx.GetUint(middleware.ContextUserID)

	resp, err := c.service.Create(userID, &req, image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

// Update — PUT /api/v1/admin/banners/:id  (multipart/form-data)
// All fields optional. Send "image" file to replace the existing one.
func (c *BannerController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.UpdateBannerRequest
	if err = ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation failed", Details: formatValidationErrors(errs),
		})
		return
	}

	image, _ := ctx.FormFile("image")

	resp, err := c.service.Update(id, &req, image)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Delete — DELETE /api/v1/admin/banners/:id
func (c *BannerController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err = c.service.Delete(id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.MessageResponse{Message: "banner deleted successfully"})
}
