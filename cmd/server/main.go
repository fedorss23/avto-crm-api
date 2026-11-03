package main

import (
	"avto-crm-api/internal/config"
	"avto-crm-api/internal/database"
	"avto-crm-api/internal/modules/auth"
	authmanagers "avto-crm-api/internal/modules/auth-managers"
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/company"
	"avto-crm-api/internal/modules/deal"
	"avto-crm-api/internal/modules/manager"
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
	stageRepo := stage.NewStageRepository(db)
	pipelineRepo := pipeline.NewPipelineRepository()
	companyRepo := company.NewCompanyRepository(db)
	managerRepo := manager.NewManagerRepository(db)

	authMConfig := authmanagers.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		MaxLoginAttempts:     5,
		LockDuration:         10 * time.Minute,
	}

	authManagerService := authmanagers.NewAuthService(managerRepo, companyRepo, jwtMaker, &authMConfig, db)
	authService := auth.NewAuthService(userRepo, jwtMaker, serviceConfig, db)
	dealService := deal.NewDealService(db, dealRepo, carRepo, pipelineRepo, stageRepo, clientRepo, managerRepo)
	carService := car.NewCarService(carRepo)
	clientSerivce := client.NewClientService(clientRepo, db)
	pipelineService := pipeline.NewPipelineService(pipelineRepo, stageRepo, db)
	userService := user.NewUserService(userRepo)
	companyService := company.NewCompanyService(companyRepo, db)

	cookieConfig := cookie.NewCookieConfig(cfg.Domain, cfg.Secure)

	authManagerHandler := authmanagers.NewAuthManagersHandler(authManagerService, cookieConfig)
	authHandler := auth.NewAuthHandler(authService, cookieConfig)
	dealHandler := deal.NewDealHandler(dealService)
	carHandler := car.NewCarHandler(carService)
	clientHandler := client.NewClientHandler(clientSerivce)
	pipelineHandler := pipeline.NewPipelineHandler(pipelineService)
	userHandler := user.NewUserHandler(userService)
	companyHandler := company.NewCompanyHandler(companyService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"time":    time.Now().Format(time.RFC3339),
			"service": "service working",
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
			aauth.Use(middleware.UserAuthMiddleware(cfg.JWTSecret))
			{
				aauth.POST("/logout", authHandler.Logout)
				aauth.GET("/profile", authHandler.GetProfile)
				aauth.POST("/change-password", authHandler.ChangePassword)
			}
		}

		authM := api.Group("/auth-manager")
		{
			authM.POST("/login", authManagerHandler.Login)
			
			authMA := api.Group("")
			authMA.Use(middleware.ManagerAuthMiddleware(cfg.JWTSecret))
			{
				authMA.POST("/refresh-tokens", authManagerHandler.RefreshTokens)
				authMA.POST("/", authManagerHandler.ChangePassword)
			}

			authMr := authM.Group("")
			authMr.Use(middleware.UserAuthMiddlewareWithCompany(cfg.JWTSecret, companyRepo.FindById))
			{
				authMr.POST("/register", authManagerHandler.Register)
			}

		}

		deal := api.Group("/deal")
		deal.Use(middleware.StandartMiddlewareWithCompany(cfg.JWTSecret, companyRepo.FindById))
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
		client.Use(middleware.StandartMiddlewareWithCompany(cfg.JWTSecret, companyRepo.FindById))
		{
			client.GET("/:clientId", clientHandler.FindById)
			client.GET("", clientHandler.FindList)
			client.PUT("/:clientId", clientHandler.Update)
		}

		pipeline := api.Group("/pipeline")
		pipeline.Use(middleware.StandartMiddlewareWithCompany(cfg.JWTSecret, companyRepo.FindById))
		{
			pipeline.GET("", pipelineHandler.FindList)
			pipeline.PUT("", pipelineHandler.Update)
			pipeline.POST("/:pipelineId/stages", pipelineHandler.AddStages)
			pipeline.DELETE("/:pipelineId/stages", pipelineHandler.DeleteStage)
		}

		users := api.Group("/users")
		users.Use(middleware.UserAuthMiddleware(cfg.JWTSecret))
		{
			users.GET("", userHandler.FindList)
			users.DELETE(":userId", userHandler.Delete)
		}

		company := api.Group("/company")
		company.Use(middleware.UserAuthMiddleware(cfg.JWTSecret))
		{
			company.POST("", companyHandler.Create)
		}
	}

	router.Run(":" + cfg.ServerPort)
}
