package handler

import (
	"net/http"
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

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// === Serve Frontend ===
	router.StaticFile("/", "./frontend/index.html") // Root → leaderboard
	router.StaticFile("/index.html", "./frontend/index.html")

	// Serve static assets (images, css, js, favicon, etc.)
	router.Static("/assets", "./frontend/assets")

	// Health check
	router.GET("/health", r.healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", r.userHandler.AddUser)
			users.GET("/:handle", r.userHandler.GetUserByHandle)
		}

		v1.GET("/leaderboard", r.userHandler.GetLeaderboard)
	}

	// Optional: SPA fallback - serve index.html for any unknown route (except API)
	router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 8 && c.Request.URL.Path[:8] == "/api/v1/" {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
			return
		}
		c.File("./frontend/index.html")
	})

	return router
}
