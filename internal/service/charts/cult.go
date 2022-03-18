package charts

import (
	"html/template"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/utils"
	"ssstatistics/pkg/chartjs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func Cult(c *gin.Context) {
	// Load data
	var cultInfo []domain.CultInfo
	r.Database.Preload(clause.Associations).Find(&cultInfo)

	// Initialize chart properties
	fields := []string{"crew_escaped", "crew_dead", "crew_total", "crew_survived", "clownabuse"}

	totalScore := utils.Slice2int32Map(fields)
	for _, score := range cultInfo {
		flatData := utils.Struct2ExpectedFieldMap(&score, fields)
		for k, v := range flatData {
			totalScore[k] += int32(v.(float64))
		}
	}

	coords := make([]chartjs.Coords, 0, len(fields))
	for k, v := range totalScore {
		coords = append(coords, chartjs.Coords{
			"x": k,
			"y": float32(v) / float32(len(cultInfo)),
		})
	}

	// Create chart
	chart := chartjs.New("bar").
		SetLabels([]string{"Crew Escaped", "Crew Dead", "Crew Total", "Crew Survived", "Clownabuse"}).
		AddDataset(chartjs.BarDataset("Crew", coords))

	c.HTML(200, "graphic.html", gin.H{
		"charts": []template.JS{chart.String()},
	})
}
