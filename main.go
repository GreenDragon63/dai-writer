package main

import (
	"dai-writer/auth"
	"dai-writer/routes"

	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("LOCAL_INSTALL") != "true" {
		auth.InitUser()
		gin.DisableConsoleColor()
		f, _ := os.Create("server.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}
	log.SetOutput(gin.DefaultWriter)
	router := gin.Default()
	router.MaxMultipartMemory = 5 * 1024 * 1024
	routes.AddPublics(router)
	routes.AddPrivates(router)
	server := os.Getenv("SERVER")
	if server != "" {
		router.Run(server)
	} else {
		router.Run("127.0.0.1:5555")
	}

}
