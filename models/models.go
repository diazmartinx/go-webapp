package models

import (
	"app/db"

	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Url        string
	Categories []Category `gorm:"foreignkey:Url;references:Url"`
}

type Category struct {
	gorm.Model
	Name  string
	Url   string
	Items []Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Item struct {
	gorm.Model
	Name       string
	Url        string
	Done       bool
	CategoryID uint
}

func MigrateCategory() {
	db.DB.AutoMigrate(&Category{})
}

func MigrateList() {
	db.DB.AutoMigrate(&List{})
}

func MigrateItem() {
	db.DB.AutoMigrate(&Item{})
}
