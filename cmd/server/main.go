package main

import (
	"log"
	"os"

	"myapp/controllers"
	"myapp/database"
	"myapp/repositories"
	"myapp/routes"
	"myapp/services"
	"myapp/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading environment variables from system")
	}

	utils.InitR2()

	db := database.Connect()
	database.AutoMigrate(db)
	database.Seed(db)

	// ── Repositories ──────────────────────────────────────
	userRepo := repositories.NewUserRepository(db)
	sportClubRepo := repositories.NewSportClubRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	permRepo := repositories.NewPermissionRepository(db)
	bannerRepo := repositories.NewBannerRepository(db)
	slotRepo := repositories.NewSlotRepository(db)

	// ── Services ──────────────────────────────────────────
	authSvc := services.NewAuthService(userRepo)
	userSvc := services.NewUserService(userRepo, sportClubRepo)
	categorySvc := services.NewCategoryService(categoryRepo)
	sportClubSvc := services.NewSportClubService(sportClubRepo, categoryRepo)
	permSvc := services.NewPermissionService(permRepo)
	bannerSvc := services.NewBannerService(bannerRepo)
	slotSvc := services.NewSlotService(slotRepo, sportClubRepo, categoryRepo)

	// ── Controllers ───────────────────────────────────────
	authCtrl := controllers.NewAuthController(authSvc)
	userCtrl := controllers.NewUserController(userSvc)
	categoryCtrl := controllers.NewCategoryController(categorySvc)
	sportClubCtrl := controllers.NewSportClubController(sportClubSvc)
	permCtrl := controllers.NewPermissionController(permSvc)
	bannerCtrl := controllers.NewBannerController(bannerSvc)
	slotCtrl := controllers.NewSlotController(slotSvc)

	// ── Router ────────────────────────────────────────────
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)

	routes.RegisterAll(r, permRepo, authCtrl, userCtrl, categoryCtrl, sportClubCtrl, slotCtrl, bannerCtrl, permCtrl)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
