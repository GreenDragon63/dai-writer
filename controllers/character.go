package controllers

import (
	"dai-writer/auth"
	"dai-writer/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong parameter"})
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
