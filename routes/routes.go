package routes

import (
	"dai-writer/auth"
	"dai-writer/controllers"

	"github.com/gin-gonic/gin"
)

func AddPublics(router *gin.Engine) {
	router.Static("/static", "static/")
	router.LoadHTMLGlob("views/*")
	router.GET("/", controllers.GetIndex)
	router.GET("/login", auth.GetLogin)
	router.POST("/login", auth.PostLogin)
	router.GET("/character/", auth.GetCurrentUser(false), controllers.Character)
	router.GET("/book/", auth.GetCurrentUser(false), controllers.Book)
	router.GET("/scene/:book", auth.GetCurrentUser(false), controllers.Scene)
	router.GET("/line/:book/:scene", auth.GetCurrentUser(false), controllers.Line)
}

func AddPrivates(router *gin.Engine) {
	privates := router.Group("/api", auth.GetCurrentUser(true))
	{
		privates.POST("/upload", controllers.UploadCharacter)
		privates.GET("/clone/:id", controllers.CloneCharacter)
		privates.GET("/character/infos/", controllers.ListCharacterInfos)
		privates.GET("/character/", controllers.ListCharacter)
		privates.GET("/character/:id", controllers.GetCharacter)
		privates.POST("/character/:id", controllers.PostCharacter)
		privates.DELETE("/character/:id", controllers.DeleteCharacter)
		privates.GET("/avatar/:id", controllers.AvatarCharacter)
		privates.GET("/export/:id", controllers.ExportBook)
		privates.GET("/book/", controllers.ListBook)
		privates.GET("/book/:id", controllers.GetBook)
		privates.POST("/book/:id", controllers.PostBook)
		privates.DELETE("/book/:id", controllers.DeleteBook)
		privates.GET("/scene/:book/", controllers.ListScene)
		privates.GET("/scene/:book/:id", controllers.GetScene)
		privates.POST("/scene/:book/:id", controllers.PostScene)
		privates.DELETE("/scene/:book/:id", controllers.DeleteScene)
		privates.GET("/line/:book/:scene/", controllers.ListLine)
		privates.GET("/line/:book/:scene/:id", controllers.GetLine)
		privates.POST("/line/:book/:scene/:id", controllers.PostLine)
		privates.DELETE("/line/:book/:scene/:id", controllers.DeleteLine)
		privates.GET("/generate/:book/:scene/:character/:id", controllers.GenerateLine)
	}
}
