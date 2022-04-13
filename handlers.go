package main

import (
	"app/db"
	"app/helpers"
	"app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func GenerateUrl() string {
	var url string

	urlExist := true
	for urlExist {
		url = helpers.GenerateRandomString(7)
		//Verify is the url is unique
		list := models.List{}
		db.DB.Where("url = ?", url).First(&list)
		if list.ID == 0 {
			urlExist = false
		}
	}

	return url
}

func CreateList(c *gin.Context) {
	url := GenerateUrl()
	var list models.List
	list.Url = url
	db.DB.Create(&list)

	// CREATE DEFAULT CATEGORY :D
	var categories = []models.Category{{Name: "🛒 Groceries", Url: url}, {Name: "🍺 Drinks", Url: url}}
	db.DB.Create(&categories)

	var histories = []models.History{{Url: url, Title: "Created: 🛒 Groceries", Changed: time.Now().UnixMilli(), TypeChange: 0},
		{Url: url, Title: "Created: 🍺 Drinks", Changed: time.Now().UnixMilli(), TypeChange: 0}}
	db.DB.Create(&histories)

	c.Redirect(http.StatusMovedPermanently, "/"+url)
	c.Abort()
}

func Home(c *gin.Context) {
	var histories []models.History

	db.DB.Order("changed desc").Limit(5).Find(&histories)
	render(c, gin.H{"histories": histories}, "home.html")
}

func ShowList(c *gin.Context) {
	url := c.Param("url")
	var list models.List
	var histories []models.History
	var deleteHistories []models.History
	var copyAlert bool // DEFAULT -> FALSE

	db.DB.Preload("Categories.Items").Where("url = ?", url).First(&list)

	if list.ID != 0 {

		db.DB.Where("url = ?", url).Order("changed desc").Limit(15).Find(&histories)

		db.DB.Where("url = ?", url).Order("changed desc").Offset(15).Find(&deleteHistories)
		db.DB.Delete(&deleteHistories) // Delete all history after 15

		if len(histories) < 10 {
			copyAlert = true
		}

		render(c, gin.H{"title": "Grocery List", "list": list, "histories": histories, "copyAlert": copyAlert}, "list.html")

	}

}

func History(c *gin.Context) {
	url := c.Param("url")
	var history models.History

	db.DB.Where("url = ?", url).Order("changed desc").First(&history)

	c.HTML(http.StatusOK, "history.html", history)
}

func Category(c *gin.Context) {

	url := c.Param("url")

	switch c.Request.Method {

	case "POST":
		{
			name := c.PostForm("name")
			newCat := models.Category{
				Name: name,
				Url:  url,
			}
			db.DB.Create(&newCat)
			db.DB.Preload("Items").Find(&newCat)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Created: '" + newCat.Name + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 0,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

			c.HTML(http.StatusOK, "category.html", newCat)
		}

	case "DELETE":
		{
			id := c.Param("id")
			var category models.Category
			db.DB.Where("id = ?", id).First(&category)
			db.DB.Unscoped().Delete(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Deleted: '" + category.Name + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 2,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

		}

	case "PUT":
		{
			id := c.Param("id")
			name := c.PostForm("name")
			var category models.Category

			db.DB.Preload("Items").Where("id = ?", id).First(&category)
			oldname := category.Name
			category.Name = name
			db.DB.Save(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Renamed: '" + oldname + "' to '" + category.Name + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 1,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

			c.HTML(http.StatusOK, "category.html", category)
		}

	}

}

func Item(c *gin.Context) {
	url := c.Param("url")

	switch c.Request.Method {
	case "POST":
		{
			u64, _ := strconv.ParseUint(c.Param("idcat"), 10, 64)
			id := uint(u64)
			name := c.PostForm("name")
			catName := c.PostForm("categoryName")
			item := models.Item{
				Name:         name,
				CategoryID:   id,
				Url:          url,
				CreatedMilis: time.Now().UnixMilli(),
			}
			db.DB.Create(&item)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Created: '" + name + "' in '" + catName + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 0,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

			db.DB.Create(&history)

			c.HTML(http.StatusOK, "item.html", item)
		}
	case "DELETE":
		{
			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
			id := uint(u64)
			var item models.Item
			db.DB.Where("id = ?", id).First(&item)
			db.DB.Unscoped().Delete(&item)

			var category models.Category
			db.DB.Where("id = ?", item.CategoryID).First(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Deleted: '" + item.Name + "' from '" + category.Name + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 2,
			}

			db.DB.Create(&history)

			// --------------------------------------------------
		}
	case "PATCH":
		{
			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
			id := uint(u64)

			var item models.Item
			db.DB.Where("id = ?", id).First(&item)
			db.DB.Unscoped().Delete(&item)

			var category models.Category
			db.DB.Where("id = ?", item.CategoryID).First(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      "Completed: '" + item.Name + "' in '" + category.Name + "'",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 3,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

			c.HTML(http.StatusOK, "history.html", history)
		}
	}
}
