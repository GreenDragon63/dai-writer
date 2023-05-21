package routes

import (
	"dai-writer/auth"
	"dai-writer/controllers"

	"github.com/gin-gonic/gin"
)

func AddPublics(router *gin.Engine) {
	router.LoadHTMLGlob("views/*")
	router.GET("/", controllers.GetIndex)
	router.POST("/login", auth.Login)
}

func AddPrivates(router *gin.Engine) {
	privates := router.Group("/private", auth.GetCurrentUser())
	{
		privates.GET("/character/:id", controllers.GetCharacter)
	}
}
