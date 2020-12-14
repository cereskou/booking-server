package models

import "gorm.io/gorm"

//CasbinRule -
type CasbinRule struct {
	gorm.Model
	PType  string `gorm:"size:40"`
	Role   string `gorm:"size:32"`
	Path   string `gorm:"size:255"`
	Method string `gorm:"size:10"` //GET PUT PUSH DELETE
}
