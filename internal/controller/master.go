package controller

import (
	"html/template"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service"
	"ssstatistics/internal/utils"
	"ssstatistics/pkg/chartjs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

const currentStatistics = "02-27-2022"

var router *gin.Engine

func runCollector(c *gin.Context) {
	startDate, _ := time.Parse("01-02-2006", currentStatistics)

	collector := service.Collector{}
	collector.CollectUrls(startDate)
	collector.CollectStatistics()

	user := c.MustGet(gin.AuthUserKey).(string)
	c.HTML(200, "secrets.html", gin.H{
		"user": user,
		"logs": collector.Logs,
	})

}

func drawChart(c *gin.Context) {
	var roots []domain.Root
	r.Database.Preload(clause.Associations).Find(&roots)

	var scores []domain.Score
	for _, v := range roots {
		if v.ServerAddress == "game.taucetistation.org:2508" {
			scores = append(scores, v.Score)
		}
	}

	labels := utils.JsonFieldNames(&scores[0], &[]string{"crew_escaped", "crew_dead", "crew_total", "crew_survived", "clownabuse"}, &utils.Ints)
	coords := make([]chartjs.Coords, 0, len(labels))

	average := utils.Slice2Map(labels)

	for _, score := range scores {
		flatData := utils.FieldNameByValue(&score)
		for k, v := range flatData {
			if !utils.Contains(labels, k) {
				continue
			}
			integer, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			average[k] += float32(integer)
		}
	}

	for k, v := range average {
		coords = append(coords, chartjs.Coords{
			"x": k,
			"y": v / float32(len(scores)),
		})
	}

	backColors := chartjs.RandomColor(len(coords))
	dataset := chartjs.Dataset{
		Label:           "Random Dataset",
		Data:            coords,
		BackgroundColor: backColors,
		BorderColor:     chartjs.RandomBorder(backColors),
		BorderWidth:     3,
	}

	charts := chartjs.New("bar").
		SetLabels(labels).
		AddDataset(&dataset).
		String()

	c.HTML(200, "graphic.html", gin.H{
		"charts": []template.JS{charts},
	})
}

func initializeRoutes() {
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	router.GET("/gamemodes", drawChart)

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
