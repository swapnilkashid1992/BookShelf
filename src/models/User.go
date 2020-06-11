package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"size:255"`
	Username     string `gorm:"type:varchar(100)”`
	Password     string `gorm:"type:varchar(100)"`
	Phone_number string `gorm:"type:varchar(10)"`
	Role         string `gorm:"type:varchar(10)"`
}
