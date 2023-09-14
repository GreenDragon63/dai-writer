package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Title",
		"prefix": os.Getenv("URL_PREFIX"),
		"js":     "index.js",
	})
}
