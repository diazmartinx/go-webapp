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

func Home(c *gin.Context) {
	render(c, gin.H{"title": "Home Page"}, "home.html")
}

func CreateList(c *gin.Context) {
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

	list := models.List{Url: url}
	db.DB.Create(&list)
	c.Redirect(http.StatusMovedPermanently, "/"+url)
	c.Abort()
}

func ShowList(c *gin.Context) {
	url := c.Param("url")
	var list models.List
	var histories []models.History

	db.DB.Preload("Categories.Items").Preload("Categories").Where("url = ?", url).Find(&list)

	db.DB.Where("url = ?", url).Order("changed desc").Limit(10).Find(&histories)

	render(c, gin.H{"title": "Grocery List", "list": list, "histories": histories}, "list.html")

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
				Title:      name + " added",
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
			db.DB.Delete(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      category.Name + " deleted",
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
			category.Name = name
			db.DB.Save(&category)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      name + " updated",
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
				Name:       name,
				CategoryID: id,
				Url:        url,
			}
			db.DB.Create(&item)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      catName + ": " + name + " added",
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
			db.DB.Delete(&item)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      item.Name + " deleted :c",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 0,
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
			item.Done = true
			db.DB.Save(&item)

			// -------------  ADD EVENT TO HISTORY ---------------

			history := models.History{
				Url:        url,
				Title:      item.Name + " done!",
				Changed:    time.Now().UnixMilli(),
				TypeChange: 3,
			}

			db.DB.Create(&history)

			// --------------------------------------------------

			c.HTML(http.StatusOK, "history.html", history)
		}
	}
}
