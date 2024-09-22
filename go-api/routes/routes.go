package routes

import (
	"gojudge/controllers"
    "github.com/gin-gonic/gin"
    "log"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/api/auth/register", controllers.Register)
    router.POST("/api/auth/login", controllers.Login)
    // router.POST("/api/submission", controllers.Submission)

    log.Println("Routes registered")
}

// StartServer initializes the server
func StartServer() {
    router := gin.Default()
    RegisterRoutes(router)

    // Start the server
    log.Fatal(router.Run(":8001")) // Listening on port 8001
}