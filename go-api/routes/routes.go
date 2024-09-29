package routes

import (
	"gojudge/controllers"
    "github.com/gin-gonic/gin"
    "log"
)

func RegisterRoutes(router *gin.Engine) {
    // USER ROUTES
	router.POST("/api/auth/register", controllers.Register)
    router.POST("/api/auth/login", controllers.Login)

    // PROBLEMS ROUTES
    router.POST("/api/problems", controllers.CreateProblem)
    router.GET("/api/problems", controllers.GetAllProblems)
    router.GET("/api/problems/:id", controllers.GetProblem) // Add this line

    // SUBMISSION ROUTES
    router.POST("/api/submission", controllers.Submission)

    
    router.GET("/api/testcase/:id", controllers.GetTestCase)

    log.Println("Routes registered")
}

// StartServer initializes the server
func StartServer() {
    router := gin.Default()
    RegisterRoutes(router)

    // Start the server
    log.Fatal(router.Run(":8001")) // Listening on port 8001
}