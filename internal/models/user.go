package models

import "time"

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"type:varchar(255);not null;unique"`
	Password string    `gorm:"type:varchar(255);not null;  "`
	Phone    string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
