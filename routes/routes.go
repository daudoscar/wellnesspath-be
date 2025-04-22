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
		user := protected.Group("/user")
		{
			user.GET("", controllers.GetUserByID)
			user.PUT("", controllers.UpdateUser)
		}

		profile := protected.Group("/profile")
		{
			profile.GET("", controllers.GetProfile)
			profile.PUT("", controllers.UpdateProfile)
		}

		exercise := protected.Group("/exercises")
		{
			exercise.GET("", controllers.GetAllExercises)
			exercise.GET("/:id", controllers.GetExerciseByID)
		}

		plan := protected.Group("/plans")
		{
			plan.POST("/generate", controllers.GenerateWorkoutPlan)
			plan.GET("", controllers.GetAllPlans)
			plan.GET("/:id", controllers.GetPlanByID)
		}
	}

	return router
}
