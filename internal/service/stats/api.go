package stats

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

func ApiMmrGET(c *gin.Context) {
	var players []*domain.Player
	r.Database.
		Select("Ckey", "MMR").
		Find(&players)

	type mmr struct {
		Ckey string
		MMR  uint
	}

	var mmrs []*mmr
	for _, player := range players {
		mmrs = append(mmrs, &mmr{
			Ckey: player.Ckey,
			MMR:  uint(player.MMR),
		})
	}

	c.JSON(200, mmrs)
}

func ApiMapsGet(c *gin.Context) {
	var MapStatistics []*domain.MapStats
	r.Database.
		Preload("MapAttributes").
		Find(&MapStatistics)

	type simpleMapStats struct {
		MapName    string
		ServerID   string
		Attributes map[string]string
	}

	var maps []*simpleMapStats
	for _, stats := range MapStatistics {
		simpleMapStat := &simpleMapStats{
			MapName:    stats.MapName,
			ServerID:   stats.ServerID,
			Attributes: make(map[string]string),
		}
		for _, attribute := range stats.MapAttributes {
			simpleMapStat.Attributes[attribute.Name] = attribute.Value
		}
		maps = append(maps, simpleMapStat)
	}

	c.JSON(200, maps)
}
