package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler   *UserHandler
	healthHandler *HealthHandler
}

func NewRouter(userHandler *UserHandler, healthHandler *HealthHandler) *Router {
	return &Router{
		userHandler:   userHandler,
		healthHandler: healthHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", r.healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", r.userHandler.AddUser)
			users.GET("/:handle", r.userHandler.GetUserByHandle)
		}

		// Leaderboard route
		v1.GET("/leaderboard", r.userHandler.GetLeaderboard)
	}

	return router
}
