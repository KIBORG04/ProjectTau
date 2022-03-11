package controller

import (
	"ssstatistics/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

const currentStatistics = "02-27-2022"

var router *gin.Engine

func runCollector() {
	startDate, _ := time.Parse("01-02-2006", currentStatistics)

	collector := service.Collector{}
	collector.CollectUrls(startDate)
	collector.CollectStatistics()
}

func initializeRoutes() {
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "SSStatistics",
		})
	})

	router.POST("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "EEEBOY",
		})
	})

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))
	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		c.HTML(200, "secrets.html", gin.H{
			"user": user,
		})
	})

}

func Run() {
	router = gin.Default()
	router.LoadHTMLGlob("../../package/templates/*")

	initializeRoutes()

	router.Run(":8080")
}
