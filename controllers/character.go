package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCharacter(c *gin.Context) {
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
	user = u.(auth.User)
	chara, err := models.LoadCharacter(&user, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}
	c.JSON(http.StatusOK, chara)
}

func PostCharacter(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool
	var character models.Character

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
	user = u.(auth.User)

	if err := c.BindJSON(&character); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ok = models.SaveJson("characters/", user.Id, id, c)
	if ok != true {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func UploadCharacter(c *gin.Context) {
	var user auth.User

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	user = u.(auth.User)

	file, err := c.FormFile("file")

	err = c.SaveUploadedFile(file, models.UploadCharacterPath(&user))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
