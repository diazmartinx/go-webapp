package main

import (
	"net/http"

	"app/models"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	models.MigrateItem()
	models.MigrateCategory()
	models.MigrateList()
	models.MigrateHistory()

	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	router.Run(":8000")

}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}

}
