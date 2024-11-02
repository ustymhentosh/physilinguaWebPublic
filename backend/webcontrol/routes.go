package webcontrol

import (
	"github.com/gin-gonic/gin"
)

// helper structs
type WhichPage struct {
	PageName string
}

type FormBody struct {
	Number  string
	Text    string
	Answer  string
	Comment string
}

// Initializes routing for web-site
func InitRouting(path_to_key string, bucketName string) {

	port := "8080"
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// setting up folder access
	r.Static("/static", "../frontend/static")
	r.LoadHTMLGlob("../frontend/templates/*")

	// initializing connection to firebase
	InitDBconnection(path_to_key, bucketName)

	// start routing for pages
	r.GET("/", WhyPageHandler)

	r.GET("/how", HowPageHandler)

	r.GET("/catalogue", CataloguePageHandler)

	r.GET("/new", NewPageHandler)

	r.POST("/formsubmit", FormSumbissionHandler)

	r.GET("/ready", ReadyProblemsHandler)

	r.GET("/new/check", PotentialyFullHandler)

	// run
	r.Run(":" + port)
}
