package models

import (
	"github.com/jinzhu/gorm"
)

// A Gallery contains image resources that are viewed by our visitors.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Title  string `gorm:"not null"`
}
