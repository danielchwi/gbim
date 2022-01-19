package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/danielchwi/gbim/backend/database"
	"github.com/danielchwi/gbim/backend/models"
)

func PersonIndex(c *gin.Context) {

	var person []models.Person

	database.DB.Find(&person)

	c.JSON(http.StatusOK, gin.H{"page": "PersonIndex", "data": person})
}

func PersonStore(c *gin.Context) {

	layout := "02-01-2006"
	birthday, err := time.Parse(layout, c.PostForm("birthday"))
	if err != nil {
		c.JSON(400, gin.H{"msg": err, "dateSupply": c.PostForm("birthday")})
		return
	}

	baptistDay, _ := time.Parse(layout, c.PostForm("baptistDay"))

	person := models.Person{
		Name:       c.PostForm("name"),
		Birthday:   birthday,
		BirthPlace: c.PostForm("birthPlace"),
		BaptistDay: baptistDay,
	}

	if err := database.DB.Create(&person).Error; err != nil {
		c.JSON(400, gin.H{"msg": "Failed to save data", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": "PersonStore", "data": person})
}

func PersonShow(c *gin.Context) {

	var person models.Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": "failed to get data"})
		return
	}
	database.DB.Find(&person)

	c.JSON(http.StatusOK, gin.H{"page": "UserShow", "value": person})
}

func PersonUpdate(c *gin.Context) {

	var person models.Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	layout := "02-01-2006"
	birthday, err := time.Parse(layout, c.PostForm("birthday"))
	if err != nil {
		c.JSON(400, gin.H{"msg": err, "dateSupply": c.PostForm("birthday")})
		return
	}

	baptistDay, _ := time.Parse(layout, c.PostForm("baptistDay"))
	if err != nil {
		c.JSON(400, gin.H{"msg": err, "dateSupply": c.PostForm("baptistDay")})
		return
	}

	person.Name = c.PostForm("name")
	person.Birthday = birthday
	person.BirthPlace = c.PostForm("birthPlace")
	person.Birthday = baptistDay

	database.DB.Updates(&person)

	c.JSON(http.StatusOK, gin.H{"page": "UserUpdate", "value": person})

}

func PersonDestroy(c *gin.Context) {

	var person models.Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	database.DB.Find(&person)
	database.DB.Delete(person)

	c.JSON(http.StatusOK, gin.H{"page": "UserDestroy", "value": fmt.Sprintf("User %s is deleted!", person.Name)})
}
