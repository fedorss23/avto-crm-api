package main

import (
	"avto-crm-api/internal/blacklist"
	"avto-crm-api/internal/config"
	"avto-crm-api/internal/database"
	"avto-crm-api/internal/modules/auth"
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/deal"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/modules/user"
	"avto-crm-api/internal/rdb"
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

	rdbClient := rdb.NewRedisClient(cfg.RedisHost + ":" + cfg.RedisPort, cfg.RedisPassword, 0)


	blacklist := blacklist.NewBlacklistService(rdbClient)

	userRepo := user.NewUserRepository(db)
	dealRepo := deal.NewDealRepository(db)
	carRepo := car.NewCarRepository(db)
	clientRepo := client.NewClientRepository(db)
	stageRepo := stage.NewStageRepository()
	pipelineRepo := pipeline.NewPipelineRepository()

	authService := auth.NewAuthService(userRepo, jwtMaker, serviceConfig, db)
	dealService := deal.NewDealService(db, dealRepo, carRepo, pipelineRepo, stageRepo, clientRepo)
	carService := car.NewCarService(carRepo)
	clientSerivce := client.NewClientService(clientRepo)
	pipelineService := pipeline.NewPipelineService(pipelineRepo, db)
	userService := user.NewUserService(userRepo)

	cookieConfig := cookie.NewCookieConfig(cfg.Domain, cfg.Secure)

	authHandler := auth.NewAuthHandler(authService, cookieConfig, blacklist)
	dealHandler := deal.NewDealHandler(dealService)
	carHandler := car.NewCarHandler(carService)
	clientHandler := client.NewClientHandler(clientSerivce)
	pipelineHandler := pipeline.NewPipelineHandler(pipelineService)
	userHandler := user.NewUserHandler(userService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"time":    time.Now().Format(time.RFC3339),
			"service": "service working",
		})
	})

	router.GET("/redis-check", func(c *gin.Context) {
		err := rdbClient.Ping(c.Request.Context())

		if err != nil {
			c.JSON(500, gin.H{
				"status": "error",
				"time": time.Now().Format(time.RFC3339),
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"time": time.Now().Format(time.RFC3339),
			"message": "pong",
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
			aauth.Use(middleware.AuthMiddleware(cfg.JWTSecret, blacklist.IsBlacklisted))
			{
				aauth.POST("/logout", authHandler.Logout)
				aauth.GET("/profile", authHandler.GetProfile)
				aauth.POST("/change-password", authHandler.ChangePassword)
			}
		}

		deal := api.Group("/deal")
		deal.Use(middleware.AuthMiddleware(cfg.JWTSecret, blacklist.IsBlacklisted))
		{
			deal.GET("", dealHandler.FindList)
			deal.GET("/total", dealHandler.GetTotalByStatus)
			deal.POST("", dealHandler.CreateFullDeal)
			deal.PUT("/:dealId", dealHandler.Update)
			deal.DELETE("/:dealId", dealHandler.Delete)
			deal.GET("/:dealId", dealHandler.FindById)
		}

		car := api.Group("/car")
		{
			car.GET("", carHandler.FindAll)
			car.POST("", carHandler.Create)
		}

		client := api.Group("/client")
		client.Use(middleware.AuthMiddleware(cfg.JWTSecret, blacklist.IsBlacklisted))
		{
			client.GET(":clientId", clientHandler.FindById)
			client.GET("/by-owner", clientHandler.FindListByOwnerId)
		}

		pipeline := api.Group("/pipeline")
		pipeline.Use(middleware.AuthMiddleware(cfg.JWTSecret, blacklist.IsBlacklisted))
		{
			pipeline.GET("", pipelineHandler.FindList)
		}

		users := api.Group("/users")
		users.Use(middleware.AdminMiddleware(cfg.JWTSecret, blacklist.IsBlacklisted))
		{
			users.GET("", userHandler.FindList)
			users.DELETE(":userId", userHandler.Delete)
		}
	}

	router.Run(":" + cfg.ServerPort)
}
