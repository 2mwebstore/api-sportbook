package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// parseID extracts and parses the :id URL param.
func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	return uint(id), err
}

// formatValidationErrors converts validator.ValidationErrors into a field→tag map.
func formatValidationErrors(err error) map[string]string {
	details := map[string]string{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			details[e.Field()] = e.Tag()
		}
	}
	return details
}
