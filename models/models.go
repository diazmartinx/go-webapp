package models

import (
	"app/db"

	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Url        string
	Categories []Category `gorm:"foreignkey:Url;references:Url"`
	//Histories  []History  `gorm:"foreignkey:Url;references:Url"` DIDNT WORK ? ;C
}

type Category struct {
	gorm.Model
	Name  string
	Url   string
	Items []Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Item struct {
	gorm.Model
	Name         string
	Url          string
	Done         bool
	CategoryID   uint
	CreatedMilis int64 // in miliseconds, then converted to javascript time in frontend
}

type History struct {
	ID         int `gorm:"primaryKey"`
	Url        string
	Title      string
	Changed    int64 // in miliseconds, then converted to javascript time in frontend
	TypeChange int   // 0(created) // 1(updated) // 2(deleted) // 3(done)
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

func MigrateHistory() {
	db.DB.AutoMigrate(&History{})
}
