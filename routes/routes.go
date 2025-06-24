package routes

import (
	"wellnesspath/controllers"
	"wellnesspath/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Global queue middleware
	router.Use(middleware.QueueMiddleware())

	// Public routes
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(middleware.AuthenticateJWT())
	{
		user := protected.Group("/user")
		{
			user.GET("", controllers.GetUserByID)
			user.PUT("", controllers.UpdateUser)
			user.DELETE("", controllers.DeleteUser)
		}

		profile := protected.Group("/profile")
		{
			profile.GET("", controllers.GetProfile)
			profile.PUT("", controllers.UpdateProfile)
			profile.DELETE("", controllers.DeleteProfile)
		}

		exercise := protected.Group("/exercises")
		{
			exercise.GET("", controllers.GetAllExercises)
			exercise.GET("/:id", controllers.GetExerciseByID)
		}

		plan := protected.Group("/plans")
		{
			plan.POST("/generate", controllers.GenerateWorkoutPlan)
			plan.GET("", controllers.GetPlanByUserID)
			plan.GET("today", controllers.GetWorkoutToday)
			plan.DELETE("", controllers.DeletePlan)
			plan.GET("/recommendations", controllers.GetRecommendedReplacements)
			plan.PUT("/replace", controllers.ReplaceExercise)
			plan.PUT("/updatereps", controllers.UpdateExerciseReps)
		}
	}

	return router
}
