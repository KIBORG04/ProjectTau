package controller

import (
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/service/parser"
	"time"

	"github.com/gin-gonic/gin"
)

const currentStatistics = "02-27-2022"

var router *gin.Engine

func runCollector(c *gin.Context) {
	startDate, _ := time.Parse("01-02-2006", currentStatistics)

	collector := parser.Collector{}
	collector.CollectUrls(startDate)
	collector.CollectStatistics()

	user := c.MustGet(gin.AuthUserKey).(string)
	c.HTML(200, "secrets.html", gin.H{
		"user": user,
		"logs": collector.Logs,
	})

}

func initializeRoutes() {
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	router.GET("/gamemodes", charts.Gamemode)
	router.GET("/cult", charts.Cult)

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))

	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.HTML(200, "secrets.html", gin.H{
			"user": user,
		})
	})

	authorized.POST("/secrets", runCollector)
}

func Run() {
	router = gin.Default()
	router.LoadHTMLGlob("../../web/templates/*")

	initializeRoutes()

	router.Run(":8080")
}
