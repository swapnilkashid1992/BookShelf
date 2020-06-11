package models

import "github.com/jinzhu/gorm"

type Booking struct {
	gorm.Model
	USERID int `gorm:"type:BigInt”`
	BOOKID int `gorm:"type:BigInt”`
}
