package controller

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/service/stats/ckey_statistics"
)

func CkeyUplinkBuysGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyUplinkBuys(c)
	c.JSON(code, result)
}
