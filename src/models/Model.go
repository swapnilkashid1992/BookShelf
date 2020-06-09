package models

import "time"

type Model struct {
	ID        uint `gorm:"primary_key";"AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
