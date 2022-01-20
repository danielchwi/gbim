package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" uri:"id" binding:"required,uuid"`
	Username string `json:"username" gorm:"unique" form:"username"`
	Password []byte `json:"-"`
	PersonId string `json:"person_id" form:"person_id"`
	Person   Person
}

func (u *User) SetPassword(password []byte) []byte {
	hashPassword, _ := bcrypt.GenerateFromPassword(password, 14)
	u.Password = hashPassword
	return hashPassword
}

func (u User) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(u.Password, password)
}

func (*User) Take(db *gorm.DB, offset int, limit int) interface{} {
	var users []User

	db.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	return users
}

func (u User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(u).Count(&total)
	return total
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
