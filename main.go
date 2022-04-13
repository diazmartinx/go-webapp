package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var router *gin.Engine

func main() {

	// IP RATE LIMITER
	rate, err := limiter.NewRateFromFormatted("4-S")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance_iprate := limiter.New(store, rate)
	middleware_iprate := mgin.NewMiddleware(instance_iprate)
	// END IP RATE LIMITER

	// DB MIGRATIONS
	//	models.MigrateItem()
	//	models.MigrateCategory()
	//	models.MigrateList()
	//	models.MigrateHistory()
	// END DB MIGRATIONS

	router = gin.Default()

	router.ForwardedByClientIP = true
	router.Use(middleware_iprate)

	//templ := template.Must(template.New("").ParseFS(embeddedFiles, "templates/*"))
	//router.SetHTMLTemplate(templ)

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	router.Run()

}

// Render one of HTML or JSON based on the 'Accept' header of the request
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
