package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/utils"
)

func Root(c *gin.Context) {
	var roots []domain.Root
	r.Database.Preload("Deaths").Find(&roots)

	var links domain.Link
	str := fmt.Sprintf("%%%d%%", roots[len(roots)-1].RoundID)
	r.Database.Where("link LIKE ?", str).First(&links)

	crewDeathsCount := make(map[string]uint, 0)
	crewDeathsSum := 0

	roleDeathsCount := make(map[string]uint, 0)
	roleDeathsSum := 0

	modesCount := make(map[string]uint, 0)
	modesSum := 0

	var alphaRoots []domain.Root
	var betaRoots []domain.Root
	var gammaRoots []domain.Root
	for _, root := range roots {
		switch root.ServerAddress {
		case ServerAlpha:
			alphaRoots = append(alphaRoots, root)
		case ServerBeta:
			betaRoots = append(betaRoots, root)
		case ServerGamma:
			gammaRoots = append(gammaRoots, root)
		}

		modesCount[root.Mode]++
		modesSum++

		for _, death := range root.Deaths {
			if slices.Contains(stationPositions, death.AssignedRole) && utils.IsDrone.FindString(death.Name) == "" {
				crewDeathsCount[death.AssignedRole]++
				crewDeathsSum++
			}
			if death.SpecialRole != "" {
				roleDeathsCount[death.SpecialRole]++
				roleDeathsSum++
			}
		}
	}

	c.HTML(200, "index.html", gin.H{
		"totalRounds":     len(roots),
		"version":         roots[len(roots)-1].Version,
		"lastDate":        links.Date,
		"alphaRounds":     len(alphaRoots),
		"betaRounds":      len(betaRoots),
		"gammaRounds":     len(gammaRoots),
		"modesCount":      modesCount,
		"modesSum":        modesSum,
		"crewDeathsCount": crewDeathsCount,
		"crewDeathsSum":   crewDeathsSum,
		"roleDeathsCount": roleDeathsCount,
		"roleDeathsSum":   roleDeathsSum,
	})
}

func Gamemode(c *gin.Context) {
	render := charts.Gamemode()

	c.HTML(200, "chart.html", gin.H{
		"charts": render,
	})
}

func Cult(c *gin.Context) {
	render := charts.Cult()

	c.HTML(200, "chart.html", gin.H{
		"charts": render,
	})
}
