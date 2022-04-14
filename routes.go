package main

func InitializeRoutes() {

	router.GET("/", Home)
	router.POST("/", CreateList)
	router.GET("/:url", ShowList)

	// CATEGORY CRUD
	router.POST("/:url/category", Category)
	router.DELETE("/:url/category/:id", Category)
	router.PUT("/:url/category/:id", Category)

	// ITEM CRUD
	router.POST("/:url/item/:idcat", Item)
	router.DELETE("/:url/item/:iditem", Item)
	router.PUT("/:url/item/:iditem", Item)   // CHANGE NAME
	router.PATCH("/:url/item/:iditem", Item) // CHANGE DONE

	// HISTORY
	router.GET("/:url/history", History)

}
