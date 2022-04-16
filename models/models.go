package models

import "app/db"

type List struct {
	ID  uint   `gorm:"primaryKey"`
	Url string `gorm:"not null; unique"`

	UpdatedAt  int
	Categories []Category `gorm:"foreignkey:Url;references:Url"`
}

type Category struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Url        string
	UpdatedAt  int
	Items      []Item      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CatChanges []CatChange `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CatChange struct {
	ID         uint `gorm:"primaryKey"`
	Title      string
	Url        string
	TypeChange int // 0(created) - 1(updated) - 2(deleted) - 3(done)

	CategoryID uint
	UpdatedAt  int
}

type Item struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Url  string

	CategoryID uint
	UpdatedAt  int
}

func Migrate() {
	db.DB.AutoMigrate(&List{}, &Category{}, &Item{}, &CatChange{})
}
