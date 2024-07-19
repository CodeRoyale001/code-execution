package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/low4ey/OJ/Golang-backend/controllers"
)

func SubmissionRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/submit", controllers.Submit())
	incomingRoutes.POST("/run", controllers.Run())
	incomingRoutes.GET("/getallsubmission", controllers.GetAllSub())
	incomingRoutes.GET("/getsubmission/question/:questionId", controllers.GetSubByQuestionId())
	incomingRoutes.GET("/getsubmission/user/:userId", controllers.GetSubByUserId())
	incomingRoutes.GET("/ping", func(c *gin.Context) {
		c.String(200, "Server is up and running :)")
	})
}
