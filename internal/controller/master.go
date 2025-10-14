package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	c "ssstatistics/internal/config"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	//"ssstatistics/internal/service/forecast" // Можно оставить для костыля
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

		base.GET("/changeling", GET(general.ChangelingGET))
		base.POST("/changeling", POST(general.ChangelingGET))

		base.GET("/finder", GET(FinderGET))

		base.GET("/player", GET(NotAPlayerGET))

		base.GET("/player/:name", GET(PlayerStatisticGET))
	}

	apiRoute := base.Group("/api")
	{
		// support query with 'ckey' parameter
		apiRoute.GET("/mmr", api.MmrGET)

		apiRoute.GET("/maps", api.MapsGET)
		apiRoute.POST("/maps", api.MapsGET)

		apiRoute.POST("/send_feedback", api.SendFeedback)

		apiRoute.GET("/heatmaps", api.HeatmapsGET)

		apiRoute.GET("/changeling", api.ChangelingGET)
		apiRoute.POST("/changeling", api.ChangelingGET)

		apiRoute.GET("/uplink", api.UplinkGET)
		apiRoute.POST("/uplink", api.UplinkGET)

		apiRoute.GET("/random_announce", api.RandomAnnounceGET)

		apiRoute.GET("/random_achievement", api.RandomAchievementGET)

		apiRoute.GET("/random_last_phrase", api.RandomLastPhraseGET)

		apiRoute.GET("/random_flavor", api.RandomFlavorGET)

		apiRoute.GET("/mode_winrates_by_month", api.ModeWinratesByMonthGET)

		apiRoute.GET("/online_stat", api.OnlineStatAvgPerDayGET)

		apiRoute.GET("/online_stat_weeks_forecast", api.ForecastHandler)

		apiRoute.GET("/online_stat_daily_forecast", api.DailyForecastHandler)

		apiRoute.GET("/one_step_forecast_history", api.OneStepForecastHistoryHandler)

		apiRoute.GET("/online_stat_max", api.OnlineStatMaxPerDayGET)

		apiRoute.GET("/online_stat_weeks", api.OnlineStatWeeksGET)

		apiRoute.GET("/online_stat_daytime", api.OnlineStatByDaytimeGET)

		apiRoute.GET("/chronicles_daytime", api.ChronicleByDatetimeGET)

		apiRoute.GET("/completion_html_by", api.CompletionHTMLByIdGET)
	}

	playerRoute := apiRoute.Group("/player")
	{
		playerRoute.GET("/uplink_buys", CkeyUplinkBuysGET)

		playerRoute.GET("/changeling_buys", CkeyChangelingBuysGET)

		playerRoute.GET("/wizard_buys", CkeyWizardBuysGET)

		playerRoute.GET("/try_find", TryFindCkeyGET)

		playerRoute.GET("/try_find_character", TryFindCharacterGET)

		playerRoute.GET("/characters", CkeyCharactersGET)

		playerRoute.GET("/ckeys_by_char", CharacterCkeysGET)

		playerRoute.GET("/roles", CkeyRolesGET)

		playerRoute.GET("/achievements", AchievementsCkeysGET)

		playerRoute.GET("/roles_rounds", AllRolesRoundsGET)

		playerRoute.GET("/mmr", CkeyMMRGET)

		playerRoute.GET("/crawler", CrawlerGET)
	}

	ss14 := apiRoute.Group("/ss14")
	{
		ss14.GET("/sponsors/:userid", SponsorsGET)
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
		authorized.POST("/add_chronicle", api.AddChronicle)
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

	err = r.Database.AutoMigrate(&domain.ForecastHistory{})
	if err != nil {
		panic(fmt.Sprintf("failed to auto-migrate database: %v", err))
	}

	initializeRoutes()

		// --- НАЧАЛО: ВРЕМЕННЫЙ КОД (КОСТЫЛЬ) ДЛЯ ПЕРВОГО РАСЧЕТА ---
	// Этот код запустит расчет прогнозов один раз при старте сервера.
	// После успешного расчета его нужно будет удалить.
	//fmt.Println("Starting one-time forecast calculation in background...")
	//go forecast.UpdateDailyForecast()
	//go forecast.UpdateWeeklyForecast()
	// --- КОНЕЦ ВРЕМЕННОГО КОДА ---

	router.Run(":8080")
}
