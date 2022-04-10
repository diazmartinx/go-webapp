package main

import (
	"app/db"
	"app/helpers"
	"app/models"
	"net/http"
	"strconv"

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
	var items []models.Item

	db.DB.Preload("Categories.Items").Preload("Categories").Where("url = ?", url).Find(&list)

	db.DB.Where("url = ?", url).Limit(10).Find(&items, "done = ?", true)

	render(c, gin.H{"title": "Grocery List", "list": list, "history": items}, "list.html")

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
			c.HTML(http.StatusOK, "category.html", newCat)
		}

	case "DELETE":
		{
			id := c.Param("id")
			var category models.Category
			db.DB.Where("id = ?", id).First(&category)
			if category.Url == url {
				db.DB.Delete(&category)
			}

		}

	case "PUT":
		{
			id := c.Param("id")
			name := c.PostForm("name")
			var category models.Category

			db.DB.Preload("Items").Where("id = ?", id).First(&category)
			category.Name = name
			db.DB.Save(&category)
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
			item := models.Item{
				Name:       name,
				CategoryID: id,
				Url:        url,
			}
			db.DB.Create(&item)

			c.HTML(http.StatusOK, "item.html", item)
		}
	case "DELETE":
		{
			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
			id := uint(u64)
			var item models.Item
			db.DB.Where("id = ?", id).First(&item)
			db.DB.Delete(&item)
		}
	case "PUT":
		{
			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
			id := uint(u64)
			name := c.PostForm("name")
			var item models.Item
			db.DB.Where("id = ?", id).First(&item)
			item.Name = name
			db.DB.Save(&item)
			c.HTML(http.StatusOK, "item.html", item)

		}

	case "PATCH":
		{
			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
			id := uint(u64)

			var item models.Item
			db.DB.Where("id = ?", id).First(&item)
			item.Done = true
			db.DB.Save(&item)
			c.HTML(http.StatusOK, "history.html", item)
		}
	}
}
