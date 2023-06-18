package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Scene(c *gin.Context) {
	c.HTML(http.StatusOK, "scene.tmpl", gin.H{
		"title": "Scene",
	})
}

func ListScene(c *gin.Context) {
	var user auth.User

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	book, err := strconv.Atoi(c.Param("book"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	scene, ok := models.ListScene(&user, book)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Scene not found"})
		return
	}
	c.JSON(http.StatusOK, scene)
}

func GetScene(c *gin.Context) {
	var user auth.User
	var id int

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	book, err := strconv.Atoi(c.Param("book"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	scene, ok := models.LoadScene(&user, book, id)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Scene not found"})
		return
	}
	c.JSON(http.StatusOK, scene)
}

func PostScene(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool
	var Scene models.Scene

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	book, err := strconv.Atoi(c.Param("book"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)

	if err := c.BindJSON(&Scene); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Scene format"})
		return
	}
	ok = models.SaveScene(&user, book, id, Scene)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Scene not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func DeleteScene(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	book, err := strconv.Atoi(c.Param("book"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad parameter"})
		return
	}
	user = u.(auth.User)
	ok = models.DeleteScene(&user, book, id)
	if ok != true {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
