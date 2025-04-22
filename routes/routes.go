package routes

import (
	"syncspend/controllers"
	"syncspend/middleware"

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
