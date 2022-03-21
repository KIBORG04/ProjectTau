package charts

import (
	"fmt"
	"html/template"
	"reflect"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/utils"
	"ssstatistics/pkg/chartjs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func makeAverageCoords[T any](fields []string, mainPool []T) []chartjs.Coords {
	totalScore := utils.Slice2Map[int32](fields)
	poolReflect := reflect.ValueOf(mainPool)
	for _, val := range mainPool {
		flatData := utils.Struct2ExpectedFieldMap(val, fields)
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
		if root.ServerAddress == "game.taucetistation.org:2506" {
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

	fmt.Println(chart)
	fmt.Println(chart.String())
	fmt.Println([]template.JS{chart.String()})

	c.HTML(200, "chart.html", gin.H{
		"charts": []template.JS{chart.String()},
	})
}
