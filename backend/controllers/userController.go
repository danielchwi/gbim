package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danielchwi/gbim/backend/database"
	"github.com/danielchwi/gbim/backend/models"
)

func UserIndex(c *gin.Context) {

	var users []models.User

	database.DB.Preload("Person").Find(&users)

	c.JSON(http.StatusOK, gin.H{"page": "UserIndex", "data": users})
}

func UserStore(c *gin.Context) {

	user := models.User{
		Username: c.PostForm("username"),
		PersonId: c.PostForm("person_id"),
	}

	if c.PostForm("password") != c.PostForm("password_confirm") {
		c.JSON(400, gin.H{"msg": "password not match"})
		return
	}

	user.SetPassword([]byte(c.PostForm("password")))

	if err := database.DB.Omit("Person").Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"msg": "Failed to save data", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": "UserStore", "data": user})
}

func UserShow(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": &user})
		return
	}
	database.DB.Preload("Person").Find(&user)

	c.JSON(http.StatusOK, gin.H{"page": "UserShow", "value": user})
}

func UserUpdate(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	user.Username = c.PostForm("username")
	user.PersonId = c.PostForm("person_id")

	if c.PostForm("password") != c.PostForm("password_confirm") {
		c.JSON(400, gin.H{"msg": "password not match"})
		return
	}

	user.SetPassword([]byte(c.PostForm("password")))

	database.DB.Omit("Person").Updates(&user)

	c.JSON(http.StatusOK, gin.H{"page": "UserUpdate", "value": user})

}

func UserDestroy(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	database.DB.Find(&user)
	database.DB.Omit("Person").Delete(user)

	c.JSON(http.StatusOK, gin.H{"page": "UserDestroy", "value": fmt.Sprintf("User %s is deleted!", user.Username)})
}
