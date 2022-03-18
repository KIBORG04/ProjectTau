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

	aspectsPool := make([]domain.Aspects, 0, len(cultInfo))
	for _, cult := range cultInfo {
		// skip cult from portal
		if cult.AnomaliesDestroyed == 0 && cult.RunesOnStation == 0 && cult.EndFavor == 1000 {
			continue
		}
		aspectsPool = append(aspectsPool, cult.Aspects)
	}
	ritesPool := make([]domain.RitenameByCount, 0, len(cultInfo))
	for _, cult := range cultInfo {
		// skip cult from portal
		if cult.AnomaliesDestroyed == 0 && cult.RunesOnStation == 0 && cult.EndFavor == 1000 {
			continue
		}
		ritesPool = append(ritesPool, cult.RitenameByCount)
	}

	// Initialize chart properties
	ritesFields := utils.JsonFieldNames(&cultInfo[1].RitenameByCount, nil)
	aspectsFields := utils.JsonFieldNames(&cultInfo[1].Aspects, nil)

	ritesCoords := makeAverageCoords(ritesFields, ritesPool)
	aspectsCoords := makeAverageCoords(aspectsFields, aspectsPool)

	// Create chart
	ritesChart := chartjs.New("bar").
		SetLabels(ritesFields).
		AddDataset(chartjs.BarDataset("Rites", ritesCoords))
	aspectChart := chartjs.New("bar").
		SetLabels(aspectsFields).
		AddDataset(chartjs.BarDataset("Aspects", aspectsCoords))

	c.HTML(200, "chart.html", gin.H{
		"charts": []template.JS{ritesChart.String(), aspectChart.String()},
	})
}
