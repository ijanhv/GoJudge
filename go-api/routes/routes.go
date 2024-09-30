package routes

import (
	"gojudge/controllers"
	"gojudge/middleware"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func RegisterRoutes(router *gin.Engine) {
	// USER ROUTES
	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)
	router.GET("/api/user/profile", middleware.CheckAuth, controllers.GetProfile)

	// PROBLEMS ROUTES
	router.POST("/api/problems", controllers.CreateProblem)
	router.GET("/api/problems", controllers.GetAllProblems)
	router.GET("/api/problems/:slug", controllers.GetProblem) // Add this line

	// SUBMISSION ROUTES
	router.POST("/api/submission", controllers.Submission)

	router.GET("/api/testcase/:id", controllers.GetTestCase)

	log.Println("Routes registered")
}

// StartServer initializes the server
func StartServer() {
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Add your frontend origin here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterRoutes(router)

	// Start the server
	log.Fatal(router.Run(":8001")) // Listening on port 8001
}
