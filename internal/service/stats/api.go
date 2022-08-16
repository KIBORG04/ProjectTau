package stats

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

func ApiMmrGET(c *gin.Context) {
	var mmrs []*domain.Player
	r.Database.
		Select("Ckey", "MMR").
		Find(&mmrs)

	c.JSON(200, mmrs)
}
