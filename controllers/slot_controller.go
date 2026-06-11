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

type SlotController struct {
	service  services.SlotService
	validate *validator.Validate
}

func NewSlotController(service services.SlotService) *SlotController {
	return &SlotController{service: service, validate: validator.New()}
}

// GetBySportClub — GET /api/v1/sport-clubs/:id/slots?page=1&limit=10&available=true
func (c *SlotController) GetBySportClub(ctx *gin.Context) {
	sportClubID, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid sport_club id"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	available := ctx.Query("available") == "true"

	resp, err := c.service.GetBySportClub(sportClubID, page, limit, available)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetByID — GET /api/v1/slots/:id
func (c *SlotController) GetByID(ctx *gin.Context) {
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

// Create — POST /api/v1/partner/slots  (multipart/form-data)
// Fields: sport_club_id, name, description, price, capacity, is_available, category_id
// File  : image (optional)
func (c *SlotController) Create(ctx *gin.Context) {
	var req dto.CreateSlotRequest
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

// Update — PUT /api/v1/partner/slots/:id  (multipart/form-data)
func (c *SlotController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.UpdateSlotRequest
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
		if err.Error() == "slot not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Delete — DELETE /api/v1/partner/slots/:id  or  /api/v1/admin/slots/:id
func (c *SlotController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err = c.service.Delete(id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "slot not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.MessageResponse{Message: "slot deleted successfully"})
}
