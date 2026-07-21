package main

import (
	"avto-crm-api/internal/config"
	"avto-crm-api/internal/database"
	"avto-crm-api/internal/modules/auth"
	"avto-crm-api/internal/modules/deal"
	"avto-crm-api/internal/modules/user"
	"avto-crm-api/pkg/cookie"
	"avto-crm-api/pkg/jwt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	jwtMaker := jwt.NewJWTMaker(cfg.JWTSecret, cfg.Issuer)

	log.Println("Connecting to database")

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatalf("Error with connecting to database")
	}

	userRepo := user.NewUserRepository(db)
	dealRepo := deal.NewDealRepository(db)
	// stageRepo := stage.NewStageRepository()
	// pipelineRepo := pipeline.NewPipelineRepository()

	serviceConfig := &auth.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		MaxLoginAttempts:     5,
		LockDuration:         10 * time.Minute,
	}

	authService := auth.NewAuthService(userRepo, jwtMaker, serviceConfig)
	dealService := deal.NewDealService(db, dealRepo)


	cookieConfig := cookie.NewCookieConfig(cfg.Domain, cfg.Secure)

	authHandler := auth.NewAuthHandler(authService, cookieConfig)
	dealHandler := deal.NewDealHandler(dealService)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"time":    time.Now().Format(time.RFC3339),
			"service": "hvjghjjgvjg",
		})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/check/:userId", authHandler.CheckAuth)
			auth.POST("/refresh-tokens", authHandler.Refresh)

			{
				auth.GET("/profile/:userId", authHandler.GetProfile)
				auth.POST("/change-password/:userId", authHandler.ChangePassword)
				auth.POST("/logout/:userId", authHandler.Logout)
			}
		}

		deal := api.Group("/deal")
		{
			deal.GET("/by-owner/:userId", dealHandler.FindDealByOwnerId)
			deal.GET("/by-client/:userId", dealHandler.FindDealByClientId)
			deal.POST("/", dealHandler.CreateFullDeal)
			deal.PUT("/", dealHandler.Update)
		}
	}
}