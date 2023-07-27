package controller

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/service/stats/api/ss14"
)

func SponsorsGET(c *gin.Context) {
	code, result := ss14.GetSponsor(c)
	c.JSON(code, result)
}
