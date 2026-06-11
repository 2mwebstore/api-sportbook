package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"myapp/dto"
	"myapp/middleware"
	"myapp/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SportClubController struct {
	service  services.SportClubService
	validate *validator.Validate
}

func NewSportClubController(service services.SportClubService) *SportClubController {
	return &SportClubController{service: service, validate: validator.New()}
}

func (c *SportClubController) Create(ctx *gin.Context) {
	var req dto.CreateSportClubRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation failed", Details: formatValidationErrors(errs)})
		return
	}
	userID := ctx.GetUint(middleware.ContextUserID)
	images := extractFiles(ctx, "images")
	resp, err := c.service.Create(userID, &req, images)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (c *SportClubController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	resp, err := c.service.GetAll(page, limit, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetByID returns full detail including available slots.
func (c *SportClubController) GetByID(ctx *gin.Context) {
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

func (c *SportClubController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var req dto.UpdateSportClubRequest
	if err = ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}
	if errs := c.validate.Struct(req); errs != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation failed", Details: formatValidationErrors(errs)})
		return
	}
	images := extractFiles(ctx, "images")
	resp, err := c.service.Update(id, &req, images)
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

func (c *SportClubController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err = c.service.Delete(id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "sport club not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.MessageResponse{Message: "sport club deleted successfully"})
}

func extractFiles(ctx *gin.Context, field string) []*multipart.FileHeader {
	form, err := ctx.MultipartForm()
	if err != nil || form == nil {
		return nil
	}
	return form.File[field]
}
