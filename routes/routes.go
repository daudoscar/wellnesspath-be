package routes

import (
	"wellnesspath/controllers"
	"wellnesspath/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	protected := router.Group("/protected")
	protected.Use(middleware.AuthenticateJWT())
	{

	}

	return router
}
