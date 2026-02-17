package router

import (
	"database/sql"

	"github.com/joezmuda/trackside-backend/internal/config"
	"github.com/joezmuda/trackside-backend/internal/handlers"
	mw "github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo, db *sql.DB, cfg *config.Config) {
	// Middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(mw.CORSMiddleware(cfg.CORSOrigins))

	// Repositories
	userRepo := repository.NewUserRepo(db)
	carRepo := repository.NewCarRepo(db)
	trackRepo := repository.NewTrackRepo(db)
	zoneRepo := repository.NewZoneRepo(db)
	reviewRepo := repository.NewReviewRepo(db)
	lapbookRepo := repository.NewLapbookRepo(db)

	// Handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	carHandler := handlers.NewCarHandler(carRepo)
	carModHandler := handlers.NewCarModHandler(carRepo)
	trackHandler := handlers.NewTrackHandler(trackRepo)
	trackImageHandler := handlers.NewTrackImageHandler(trackRepo)
	trackReviewHandler := handlers.NewTrackReviewHandler(trackRepo, reviewRepo)
	trackZoneHandler := handlers.NewTrackZoneHandler(trackRepo, zoneRepo)
	zoneTipHandler := handlers.NewZoneTipHandler(zoneRepo)
	lapbookHandler := handlers.NewLapbookHandler(lapbookRepo, carRepo, trackRepo)
	profileHandler := handlers.NewProfileHandler(userRepo)
	uploadHandler := handlers.NewUploadHandler(cfg.UploadDir)
	adminHandler := handlers.NewAdminHandler(trackRepo, userRepo, cfg.DataDir)

	// Auth middleware
	authMW := mw.AuthMiddleware(cfg.JWTSecret)

	// ─── Public routes ──────────────────────────────────────────────────────────
	api := e.Group("/api")

	// Auth
	api.POST("/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Public track endpoints
	api.GET("/tracks", trackHandler.List)
	api.GET("/tracks/:id", trackHandler.GetByID)
	api.GET("/tracks/:id/images", trackImageHandler.List)

	// ─── Protected routes ───────────────────────────────────────────────────────
	auth := api.Group("", authMW)

	// Cars
	auth.GET("/cars", carHandler.List)
	auth.POST("/cars", carHandler.Create)
	auth.PUT("/cars/:id", carHandler.Update)
	auth.DELETE("/cars/:id", carHandler.Delete)

	// Car mods
	auth.POST("/cars/:id/mods", carModHandler.Create)
	auth.DELETE("/cars/:id/mods/:modId", carModHandler.Delete)

	// Tracks (protected)
	auth.POST("/tracks", trackHandler.Create)
	auth.PATCH("/tracks/:id", trackHandler.Update)

	// Track images (protected)
	auth.POST("/tracks/:id/images", trackImageHandler.Create)
	auth.DELETE("/tracks/:id/images", trackImageHandler.Delete)

	// Track reviews
	auth.POST("/tracks/:id/reviews", trackReviewHandler.Create)

	// Track zones
	auth.POST("/tracks/:id/zones", trackZoneHandler.Create)
	auth.PATCH("/tracks/:id/zones/:zoneId", trackZoneHandler.Update)
	auth.DELETE("/tracks/:id/zones/:zoneId", trackZoneHandler.Delete)

	// Zone tips
	auth.POST("/tracks/:id/zones/:zoneId/tips", zoneTipHandler.Create)

	// Lapbook
	auth.GET("/lapbook", lapbookHandler.List)
	auth.POST("/lapbook", lapbookHandler.Create)
	auth.DELETE("/lapbook/:id", lapbookHandler.Delete)

	// Profile
	auth.GET("/profile", profileHandler.Get)
	auth.PUT("/profile", profileHandler.Update)

	// Upload
	auth.POST("/upload", uploadHandler.Upload)

	// Admin
	auth.POST("/admin/sync-tracks", adminHandler.SyncTracks)

	// ─── Static file serving ────────────────────────────────────────────────────
	e.GET("/uploads/*", handlers.ServeUploads(cfg.UploadDir))
}
