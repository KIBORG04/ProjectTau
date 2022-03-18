package charts

import (
	"html/template"
	"reflect"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/utils"
	"ssstatistics/pkg/chartjs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func makeAverageCoords(fields []string, mainPool interface{}) []chartjs.Coords {
	if !utils.IsSlice(mainPool) {
		return nil
	}

	totalScore := utils.Slice2int32Map(fields)
	poolReflect := reflect.ValueOf(mainPool)
	for i := 0; i < poolReflect.Len(); i++ {
		flatData := utils.Struct2ExpectedFieldMap(poolReflect.Index(i).Interface(), fields)
		for k, v := range flatData {
			totalScore[k] += int32(v.(float64))
		}
	}

	coords := make([]chartjs.Coords, 0, len(fields))
	for k, v := range totalScore {
		coords = append(coords, chartjs.Coords{
			"x": k,
			"y": float32(v) / float32(poolReflect.Len()),
		})
	}

	return coords
}

func Gamemode(c *gin.Context) {
	// Load data
	var roots []domain.Root
	r.Database.Preload(clause.Associations).Find(&roots)

	var scores []domain.Score
	for _, root := range roots {
		if root.ServerAddress == "game.taucetistation.org:2508" {
			scores = append(scores, root.Score)
		}
	}

	// Initialize chart properties
	fields := []string{"crew_escaped", "crew_dead", "crew_total", "crew_survived", "clownabuse"}

	coords := makeAverageCoords(fields, scores)

	// Create chart
	chart := chartjs.New("bar").
		SetLabels(fields).
		AddDataset(chartjs.BarDataset("Crew", coords))

	c.HTML(200, "graphic.html", gin.H{
		"charts": []template.JS{chart.String()},
	})
}
