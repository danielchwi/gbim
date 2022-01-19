package database

import (
	"github.com/danielchwi/gbim/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(sqlite.Open("file:database/gorm.db"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to databases")
	}

	DB = database

	database.AutoMigrate(
		models.User{},
		models.Person{},
	)
}
