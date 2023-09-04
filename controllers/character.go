package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Character(c *gin.Context) {
	c.HTML(http.StatusOK, "character.tmpl", gin.H{
		"title": "Character",
	})
}

func ListCharacterInfos(c *gin.Context) {
	var user auth.User

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	user = u.(auth.User)
	chara, ok := models.ListCharacterInfos(&user)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}
	c.JSON(http.StatusOK, chara)
}

func ListCharacter(c *gin.Context) {
	var user auth.User

	u, ok := c.Get("current_user")
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	user = u.(auth.User)
	chara, ok := models.ListCharacter(&user)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}
	c.JSON(http.StatusOK, chara)
}

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
	chara, ok := models.LoadCharacter(&user, id)
	if ok != true {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character format"})
		return
	}

	ok = models.SaveCharacter(&user, id, character)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func DeleteCharacter(c *gin.Context) {
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
	user = u.(auth.User)
	ok = models.DeleteCharacter(&user, id)
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
	pngFile := models.UploadCharacterPath(&user)
	err = c.SaveUploadedFile(file, pngFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	ok = models.DecodeCharacter(pngFile)
	if ok != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not a png character card"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func CloneCharacter(c *gin.Context) {
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
	chara, ok := models.LoadCharacter(&user, id)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}
	chara.Name = chara.Name + "|cloned"
	ok = models.SaveCharacter(&user, 0, *chara)
	if ok != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cannot save character"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func AvatarCharacter(c *gin.Context) {
	var user auth.User
	var id int
	var ok bool
	var pngFile, sfw string

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
	sfw = os.Getenv("SFW")
	if sfw == "true" {
		c.File("static/img/placeholder.svg")
		return
	}
	user = u.(auth.User)
	pngFile = models.AvatarCharacterPath(&user, id)
	if _, err := os.Stat(pngFile); err == nil {
		log.Println(pngFile)
		c.File(pngFile)
	} else {
		c.File("static/img/placeholder.svg")
	}
}
