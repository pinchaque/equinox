package routers

import (
	"equinox/internal/ctl"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	public := router.Group("/")
	{
		public.GET("/ping", ctl.Ping)
	}

	// Protected routes
	protected := router.Group("/api")
	// protected.Use(middleware.AuthMiddleware()) TODO add auth
	{
		protected.POST("/points", ctl.PointAdd)
	}

	return router
}
