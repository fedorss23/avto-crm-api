package main

import (
	"avto-crm-api/internal/config"
	"avto-crm-api/internal/database"
	"avto-crm-api/internal/modules/auth"
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/deal"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/modules/user"
	"avto-crm-api/pkg/cookie"
	"avto-crm-api/pkg/jwt"
	"avto-crm-api/pkg/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cfg := config.LoadConfig()

	serviceConfig := &auth.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		MaxLoginAttempts:     5,
		LockDuration:         10 * time.Minute,
	}

	if cfg.Version == "dev" {
		serviceConfig.AccessTokenDuration = 60 * 24 * time.Minute
	}

	jwtMaker := jwt.NewJWTMaker(cfg.JWTSecret, cfg.Issuer, serviceConfig.AccessTokenDuration, serviceConfig.RefreshTokenDuration)

	log.Println("Connecting to database")

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatalf("Error with connecting to database")
	}

	userRepo := user.NewUserRepository(db)
	dealRepo := deal.NewDealRepository(db)
	carRepo := car.NewCarRepository(db)
	clientRepo := client.NewClientRepository(db)
	stageRepo := stage.NewStageRepository()
	pipelineRepo := pipeline.NewPipelineRepository()

	authService := auth.NewAuthService(userRepo, jwtMaker, serviceConfig)
	dealService := deal.NewDealService(db, dealRepo, carRepo, pipelineRepo, stageRepo, clientRepo)
	carService := car.NewCarService(carRepo)
	clientSerivce := client.NewClientService(clientRepo)
	pipelineService := pipeline.NewPipelineService(pipelineRepo, db)
	userService := user.NewUserService(userRepo)

	cookieConfig := cookie.NewCookieConfig(cfg.Domain, cfg.Secure)

	authHandler := auth.NewAuthHandler(authService, cookieConfig)
	dealHandler := deal.NewDealHandler(dealService)
	carHandler := car.NewCarHandler(carService)
	clientHandler := client.NewClientHandler(clientSerivce)
	pipelineHandler := pipeline.NewPipelineHandler(pipelineService)
	userHandler := user.NewUserHandler(userService)

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

			aauth := auth.Group("")
			aauth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				aauth.GET("/profile", authHandler.GetProfile)
				aauth.POST("/change-password", authHandler.ChangePassword)
				aauth.POST("/logout", authHandler.Logout)
			}
		}

		deal := api.Group("/deal")
		deal.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			deal.GET("", dealHandler.FindAll)
			deal.GET("/by-owner", dealHandler.FindDealByOwnerId)
			deal.GET("/by-client/:clientId", dealHandler.FindDealByClientId)
			deal.POST("", dealHandler.CreateFullDeal)
			deal.PUT("", dealHandler.Update)
			deal.POST("/:dealId/next-stage", dealHandler.SetNextStage)
			deal.DELETE("/:dealId", dealHandler.Delete)
			deal.GET("/by-id/:dealId", dealHandler.FindById)
			deal.POST("/cancel/:dealId", dealHandler.CancelDeal)
			deal.POST("/avtivate/:dealId", dealHandler.ActiveDeal)
			deal.POST("/change-stage/:dealId", dealHandler.ChangeStage)
		}

		car := api.Group("/car")
		{
			car.GET("", carHandler.FindAll)
			car.POST("", carHandler.Create)
		}

		client := api.Group("/client")
		client.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			client.GET(":clientId", clientHandler.FindById)
		}

		pipeline := api.Group("/pipeline")
		pipeline.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			pipeline.GET("", pipelineHandler.FindList)
		}

		users := api.Group("/users")
		users.Use(middleware.AdminMiddleware(cfg.JWTSecret))
		{
			users.GET("", userHandler.FindList)
			users.DELETE(":userId", userHandler.Delete)
		}
	}

	router.Run(":" + cfg.ServerPort)
}
