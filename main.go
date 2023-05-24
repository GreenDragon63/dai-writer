package main

import (
	"dai-writer/routes"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("LOCAL_INSTALL") != "true" {
		gin.DisableConsoleColor()
		f, _ := os.Create("log/server.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}
	log.SetOutput(gin.DefaultWriter)
	router := gin.Default()
	router.MaxMultipartMemory = 5 * 1024 * 1024
	routes.AddPublics(router)
	routes.AddPrivates(router)
	router.Run(":5555")
}
