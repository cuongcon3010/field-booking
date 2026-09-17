package main

import (
	"log"

	"field-booking/backend/internal/config"
	"field-booking/backend/internal/database"
	"field-booking/backend/internal/handler"
	"field-booking/backend/internal/middleware"
	"field-booking/backend/internal/repository"
	"field-booking/backend/internal/service"
	"field-booking/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Tải cấu hình môi trường
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Kiểm tra cấu hình Google OAuth khi khởi động
	if err := cfg.ValidateOAuth(); err != nil {
		log.Printf("[WARNING] Google OAuth configuration is missing: %v. Google authentication endpoints will be disabled.", err)
	}

	// 3. Kết nối Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 4. Khởi tạo tầng Repository, Service, Handler
	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationHours)
	oauthProvider := service.NewGoogleOAuthProvider(cfg)
	authService := service.NewAuthService(userRepo, oauthRepo, jwtService, db)
	authHandler := handler.NewAuthHandler(authService, oauthProvider, cfg)

	// 5. Cấu hình Gin Router
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Đăng ký CORS Middleware
	r.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "field-booking-backend",
		})
	})

	// API Authentication routes
	apiAuth := r.Group("/api/auth")
	{
		apiAuth.GET("/google", authHandler.GoogleLogin)
		apiAuth.GET("/google/callback", authHandler.GoogleCallback)
		apiAuth.POST("/logout", authHandler.Logout)

		// Endpoint yêu cầu đăng nhập
		apiAuth.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.GetMe)
	}

	// 6. Khởi chạy HTTP server
	addr := ":" + cfg.Port
	log.Printf("Starting backend server on port %s (environment: %s)...", cfg.Port, cfg.Environment)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server stopped unexpectedly: %v", err)
	}
}
