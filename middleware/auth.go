package middleware

import (
	"net/http"
	"strings"

	"myapp/dto"
	"myapp/repositories"
	"myapp/utils"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID = "userID"
	ContextRole   = "role"
	ContextEmail  = "email"
	ContextPhone  = "phone"
)

// AuthRequired validates the Bearer token and stores claims in context.
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "authorization header missing",
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "authorization header format must be: Bearer <token>",
			})
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "invalid or expired token",
			})
			return
		}

		ctx.Set(ContextUserID, claims.UserID)
		ctx.Set(ContextRole, claims.Role)
		ctx.Set(ContextEmail, claims.Email)
		ctx.Set(ContextPhone, claims.Phone)
		ctx.Next()
	}
}

// RoleRequired checks the caller's role against a static list of allowed roles.
// Use this for simple role-based gates. Must be used after AuthRequired().
func RoleRequired(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(ctx *gin.Context) {
		role := ctx.GetString(ContextRole)
		if !allowed[role] {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "you do not have permission to access this resource",
			})
			return
		}
		ctx.Next()
	}
}

// PermissionRequired checks dynamically from the DB whether the caller's role
// has the given permission. Must be used after AuthRequired().
func PermissionRequired(permRepo repositories.PermissionRepository, permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString(ContextRole)
		if role == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "authorization required",
			})
			return
		}

		ok, err := permRepo.HasPermission(role, permission)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "permission check failed",
			})
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "permission denied: " + permission,
			})
			return
		}
		ctx.Next()
	}
}
