package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	ID         string    `json:"id" uri:"id" binding:"uuid"`
	Name       string    `json:"name"`
	Birthday   time.Time `json:"birthday"`
	BirthPlace string    `json:"birth_place"`
	BaptistDay time.Time `json:"baptist_day"`
}

func (*Person) Take(db *gorm.DB, offset int, limit int) interface{} {
	var person []Person

	db.Preload("Role").Offset(offset).Limit(limit).Find(&person)

	return person
}

func (p Person) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(p).Count(&total)
	return total
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
