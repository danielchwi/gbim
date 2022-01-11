package controllers

import (
	"net/http"

	"github.com/danielchwi/gbim/backend/database"
	"github.com/danielchwi/gbim/backend/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	user := models.User{
		Username: c.PostForm("username"),
	}

	if c.PostForm("password") != c.PostForm("password_confirm") {
		c.JSON(400, gin.H{"msg": "password not match"})
		return
	}

	user.SetPassword([]byte(c.PostForm("password")))

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"msg": "Failed to save data", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": "UserStore", "data": user})
}

func Login(c *gin.Context) {

	user := models.User{
		Username: c.PostForm("username"),
	}

	if err := database.DB.Find(&user).Error; err != nil {
		c.JSON(400, gin.H{"msg": "User Not Found", "err": err})
		return
	}

	if err := user.ComparePassword([]byte(c.PostForm("password"))); err != nil {
		c.JSON(400, gin.H{"msg": "Password not match"})
		return
	}

	c.SetCookie("cookieName", "name", 10, "/", "localhost", true, true)
}
