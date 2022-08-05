package controller

import (
	"github.com/gin-gonic/gin"
	c "ssstatistics/internal/config"
	"ssstatistics/internal/service/cleaning"
	"ssstatistics/internal/service/parser"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/service/tops"
)

var router *gin.Engine

func runUpdateDB(c *gin.Context) {
	var logs []string

	for _, callback := range RegularCallbacks {
		callbackLogs := callback()
		for _, s := range callbackLogs {
			logs = append(logs, s)
		}
	}

	user := c.MustGet(gin.AuthUserKey).(string)
	c.HTML(200, "secrets.html", gin.H{
		"user": user,
		"logs": logs,
	})
}

func initializeRoutes() {
	base := router.Group("")

	GET := func(f func(*gin.Context) (int, string, gin.H)) func(*gin.Context) {
		return func(c *gin.Context) {
			stats.BasicGET(c, f)
		}
	}

	POST := func(f func(*gin.Context) (int, string, gin.H)) func(*gin.Context) {
		return func(c *gin.Context) {
			stats.BasicPOST(c, f)
		}
	}

	{
		base.GET("/", GET(stats.RootGET))
		base.POST("/", POST(stats.RootGET))

		base.GET("/gamemodes", GET(stats.GamemodesGET))
		base.POST("/gamemodes", POST(stats.GamemodesGET))

		base.GET("/uplink", GET(stats.UplinkGET))
		base.POST("/uplink", POST(stats.UplinkGET))

		base.GET("/objectives", GET(stats.ObjectivesGET))
		base.POST("/objectives", POST(stats.ObjectivesGET))

		base.GET("/rounds", GET(stats.RoundsGET))
		base.POST("/rounds", POST(stats.RoundsGET))

		base.GET("/round/:id", GET(stats.RoundGET))
		base.GET("/round", GET(stats.RoundsGET))

		base.GET("/tops", GET(stats.TopsGET))
		base.POST("/tops", POST(stats.TopsGET))

		base.GET("/cult", GET(stats.Cult))
	}

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := base.Group("/admin", gin.BasicAuth(gin.Accounts{
		c.Config.AdminConfig.Login: c.Config.AdminConfig.Password,
	}))
	{
		// hit "localhost:8080/admin/secrets
		authorized.GET("/secrets", func(c *gin.Context) {
			user := c.MustGet(gin.AuthUserKey).(string)
			c.HTML(200, "secrets.html", gin.H{
				"user": user,
			})
		})

		authorized.POST("/secrets", runUpdateDB)
	}

}

func Run() {
	router = gin.Default()
	err := router.SetTrustedProxies([]string{c.Config.Proxy})
	if err != nil {
		panic(err)
	}
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/web/static/", "./web/static")

	initializeRoutes()

	router.Run(":8080")
}

type RegularCallback func() []string

var RegularCallbacks []RegularCallback

func InitializeRegularCallbacks() {
	RegularCallbacks = append(RegularCallbacks, parser.RunRoundCollector)
	RegularCallbacks = append(RegularCallbacks, tops.ParseTopData)
	RegularCallbacks = append(RegularCallbacks, cleaning.CleanAnnounces)
}
