package routes

import (
	"myapp/controllers"
	"myapp/middleware"
	"myapp/models"
	"myapp/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterAll(
	r *gin.Engine,
	permRepo repositories.PermissionRepository,
	authCtrl *controllers.AuthController,
	userCtrl *controllers.UserController,
	catCtrl *controllers.CategoryController,
	scCtrl *controllers.SportClubController,
	slotCtrl *controllers.SlotController,
	bannerCtrl *controllers.BannerController,
	permCtrl *controllers.PermissionController,
) {
	auth := middleware.AuthRequired()
	client := middleware.RoleRequired(string(models.RoleClient), string(models.RolePartner))
	partner := middleware.RoleRequired(string(models.RolePartner))
	admin := middleware.RoleRequired(string(models.RoleAdmin))

	perm := func(p string) gin.HandlerFunc {
		return middleware.PermissionRequired(permRepo, p)
	}

	// ═══════════════════════════════════════════════════════
	// AUTH
	// ═══════════════════════════════════════════════════════
	r.POST("/api/v1/auth/register", authCtrl.Register)
	r.POST("/api/v1/auth/login", authCtrl.Login)
	r.POST("/api/v1/admin/auth/login", authCtrl.Login)

	// ═══════════════════════════════════════════════════════
	// USERS
	// ═══════════════════════════════════════════════════════
	cp := r.Group("/api/v1/users", auth, client)
	{
		cp.GET("/me", perm(models.PermUserReadSelf), userCtrl.GetProfile)
		cp.PUT("/me", perm(models.PermUserUpdateSelf), userCtrl.UpdateProfile)
		cp.POST("/me/avatar", perm(models.PermUserUpdateSelf), userCtrl.UploadAvatar)
		cp.GET("/me/favorites", perm(models.PermUserReadSelf), userCtrl.GetFavorites)
		cp.POST("/me/favorites/:sport_club_id", perm(models.PermUserUpdateSelf), userCtrl.AddFavorite)
		cp.DELETE("/me/favorites/:sport_club_id", perm(models.PermUserUpdateSelf), userCtrl.RemoveFavorite)
	}
	au := r.Group("/api/v1/admin/users", auth, admin)
	{
		au.GET("/me", perm(models.PermUserReadSelf), userCtrl.GetProfile)
		au.PUT("/me", perm(models.PermUserUpdateSelf), userCtrl.UpdateProfile)
		au.POST("/me/avatar", perm(models.PermUserUpdateSelf), userCtrl.UploadAvatar)
	}

	// ═══════════════════════════════════════════════════════
	// CATEGORIES
	// ═══════════════════════════════════════════════════════
	r.GET("/api/v1/categories", catCtrl.GetAll)
	r.GET("/api/v1/categories/:id", catCtrl.GetByID)

	ac := r.Group("/api/v1/admin/categories", auth, admin)
	{
		ac.GET("", perm(models.PermCategoryRead), catCtrl.GetAll)
		ac.GET("/:id", perm(models.PermCategoryRead), catCtrl.GetByID)
		ac.POST("", perm(models.PermCategoryCreate), catCtrl.Create)
		ac.PUT("/:id", perm(models.PermCategoryUpdate), catCtrl.Update)
		ac.DELETE("/:id", perm(models.PermCategoryDelete), catCtrl.Delete)
	}

	// ═══════════════════════════════════════════════════════
	// SPORT CLUBS
	// ═══════════════════════════════════════════════════════

	// GET list — no slots (light)
	r.GET("/api/v1/sport-clubs", scCtrl.GetAll)
	// GET by ID — includes available slots
	r.GET("/api/v1/sport-clubs/:id", scCtrl.GetByID)

	ps := r.Group("/api/v1/partner/sport-clubs", auth, partner)
	{
		ps.POST("", perm(models.PermSportClubCreate), scCtrl.Create)
		ps.PUT("/:id", perm(models.PermSportClubUpdate), scCtrl.Update)
		ps.DELETE("/:id", perm(models.PermSportClubDelete), scCtrl.Delete)
	}
	as := r.Group("/api/v1/admin/sport-clubs", auth, admin)
	{
		as.GET("", perm(models.PermSportClubRead), scCtrl.GetAll)
		as.GET("/:id", perm(models.PermSportClubRead), scCtrl.GetByID)
		as.POST("", perm(models.PermSportClubCreate), scCtrl.Create)
		as.PUT("/:id", perm(models.PermSportClubUpdate), scCtrl.Update)
		as.DELETE("/:id", perm(models.PermSportClubDelete), scCtrl.Delete)
	}

	// ═══════════════════════════════════════════════════════
	// SLOTS
	// ═══════════════════════════════════════════════════════

	// Public: list slots for a sport club  ?page=1&limit=10&available=true
	r.GET("/api/v1/sport-clubs/:id/slots", slotCtrl.GetBySportClub)
	// Public: single slot detail
	r.GET("/api/v1/slots/:id", slotCtrl.GetByID)

	// Partner: create / update / delete  (sport_club_id in request body)
	partnerSlot := r.Group("/api/v1/partner/slots", auth, partner)
	{
		partnerSlot.POST("", perm(models.PermSlotCreate), slotCtrl.Create)
		partnerSlot.PUT("/:id", perm(models.PermSlotUpdate), slotCtrl.Update)
		partnerSlot.DELETE("/:id", perm(models.PermSlotDelete), slotCtrl.Delete)
	}

	// Admin: full CRUD  (sport_club_id in request body)
	adminSlot := r.Group("/api/v1/admin/slots", auth, admin)
	{
		adminSlot.POST("", perm(models.PermSlotCreate), slotCtrl.Create)
		adminSlot.GET("/:id", perm(models.PermSlotRead), slotCtrl.GetByID)
		adminSlot.PUT("/:id", perm(models.PermSlotUpdate), slotCtrl.Update)
		adminSlot.DELETE("/:id", perm(models.PermSlotDelete), slotCtrl.Delete)
	}
	adminSlotList := r.Group("/api/v1/admin/sport-clubs/:id/slots", auth, admin)
	{
		adminSlotList.GET("", perm(models.PermSlotRead), slotCtrl.GetBySportClub)
	}

	// ═══════════════════════════════════════════════════════
	// BANNERS
	// ═══════════════════════════════════════════════════════
	r.GET("/api/v1/banners/active", bannerCtrl.GetActive)

	ab := r.Group("/api/v1/admin/banners", auth, admin)
	{
		ab.GET("", perm(models.PermBannerRead), bannerCtrl.GetAll)
		ab.GET("/:id", perm(models.PermBannerRead), bannerCtrl.GetByID)
		ab.POST("", perm(models.PermBannerCreate), bannerCtrl.Create)
		ab.PUT("/:id", perm(models.PermBannerUpdate), bannerCtrl.Update)
		ab.DELETE("/:id", perm(models.PermBannerDelete), bannerCtrl.Delete)
	}

	// ═══════════════════════════════════════════════════════
	// PERMISSIONS
	// ═══════════════════════════════════════════════════════
	ap := r.Group("/api/v1/admin/permissions", auth, admin)
	{
		ap.GET("", permCtrl.GetAllPermissions)
		ap.GET("/roles/:role", permCtrl.GetRolePermissions)
		ap.PUT("/roles/:role", permCtrl.SetRolePermissions)
	}
}
