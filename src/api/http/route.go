package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/domain/entities"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(
	router *gin.Engine, // FIX: Accept existing router as parameter
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtService *auth.JWTService,
) *gin.Engine {
	// FIX: Use the provided router instead of creating a new one
	// router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Setup Swagger
	router.GET("/swaggers/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", authHandler.GetProfile)
				users.PUT("/profile", authHandler.UpdateProfile)
				users.POST("/change-password", authHandler.ChangePassword)

				// Admin-only routes
				admin := users.Group("")
				admin.Use(middleware.RoleMiddleware(entities.RoleAdmin))
				{
					admin.GET("", userHandler.ListUsers)
					admin.GET("/:id", userHandler.GetUser)
				}
			}
		}
	}

	return router
}
