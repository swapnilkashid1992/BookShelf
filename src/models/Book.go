package models

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	BookName    string `gorm:"type:varchar(100)”`
	Auther_Name string `gorm:"type:varchar(100)”`
	IsAvailable bool   `gorm:"type:boolean"`
}
