package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id       uint   `json:"id" uri:"id" binding:"required"`
	Username string `json:"username" gorm:"unique" form:"Username"`
	Password []byte `json:"-"`
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