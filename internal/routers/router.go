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
	/* future expansion
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
	  protected.GET("/users", controllers.GetUsers)
	  protected.POST("/users", controllers.CreateUser)
	  protected.GET("/users/:id", controllers.GetUserByID)
	  protected.PUT("/users/:id", controllers.UpdateUser)
	  protected.DELETE("/users/:id", controllers.DeleteUser)
	}
	*/
	return router
}
