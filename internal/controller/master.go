package controller

import (
	"github.com/gin-gonic/gin"
	c "ssstatistics/internal/config"
	"ssstatistics/internal/service/stats/api"
	"ssstatistics/internal/service/stats/general"
)

var router *gin.Engine

func runUpdateDB(c *gin.Context) {
	logs := StartUpdaters()

	user := c.MustGet(gin.AuthUserKey).(string)
	c.HTML(200, "secrets.html", gin.H{
		"user": user,
		"logs": logs,
	})
}

func initializeRoutes() {
	GET := func(f func(*gin.Context) (int, string, gin.H)) func(*gin.Context) {
		return func(c *gin.Context) {
			general.BasicGET(c, f)
		}
	}

	POST := func(f func(*gin.Context) (int, string, gin.H)) func(*gin.Context) {
		return func(c *gin.Context) {
			general.BasicPOST(c, f)
		}
	}

	base := router.Group("")
	{
		base.GET("/", GET(general.RootGET))
		base.POST("/", POST(general.RootGET))

		base.GET("/gamemodes", GET(general.GamemodesGET))
		base.POST("/gamemodes", POST(general.GamemodesGET))

		base.GET("/uplink", GET(general.UplinkGET))
		base.POST("/uplink", POST(general.UplinkGET))

		base.GET("/objectives", GET(general.ObjectivesGET))
		base.POST("/objectives", POST(general.ObjectivesGET))

		base.GET("/rounds", GET(general.RoundsGET))
		base.POST("/rounds", POST(general.RoundsGET))

		base.GET("/round/:id", GET(general.RoundGET))
		base.GET("/round", GET(general.RoundsGET))

		base.GET("/tops", GET(general.TopsGET))
		base.POST("/tops", POST(general.TopsGET))

		base.GET("/mmr", GET(general.MmrGET))
		base.POST("/mmr", POST(general.MmrGET))

		base.GET("/maps", GET(general.MapsGET))
		base.POST("/maps", POST(general.MapsGET))

		base.GET("/cult", GET(general.Cult))

		base.GET("/feedback", GET(general.FeedbackGET))

		base.GET("/heatmaps", GET(general.HeatmapsGET))

		base.GET("/changling", GET(general.ChanglingGET))
		base.POST("/changling", POST(general.ChanglingGET))
	}

	apiRoute := base.Group("/api")
	{
		// support query with 'ckey' parameter
		apiRoute.GET("/mmr", api.MmrGET)

		apiRoute.GET("/maps", api.MapsGET)
		apiRoute.POST("/maps", api.MapsGET)

		apiRoute.POST("/send_feedback", api.SendFeedback)

		apiRoute.GET("/heatmaps", api.HeatmapsGET)

		apiRoute.GET("/changling", api.ChanglingGET)
		apiRoute.POST("/changling", api.ChanglingGET)

		apiRoute.GET("/uplink", api.UplinkGET)
		apiRoute.POST("/uplink", api.UplinkGET)

		apiRoute.GET("/random_announce", api.RandomAnnounceGET)

		apiRoute.GET("/random_achievement", api.RandomAchievementGET)

		apiRoute.GET("/random_last_phrase", api.RandomLastPhraseGET)

		apiRoute.GET("/random_flavor", api.RandomFlavorGET)

		apiRoute.GET("/uplink_buys", CkeyUplinkBuysGET)

		apiRoute.GET("/mode_winrates_by_month", api.ModeWinratesByMonthGET)
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
