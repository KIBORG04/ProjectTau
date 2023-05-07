package controller

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/service/stats/ckey_statistics"
)

func FinderGET(c *gin.Context) (int, string, gin.H) {
	return 200, "finder.html", gin.H{}
}

func CkeyUplinkBuysGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyUplinkBuys(c)
	c.JSON(code, result)
}

func CkeyChanglingBuysGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyChanglingBuys(c)
	c.JSON(code, result)
}

func CkeyWizardBuysGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyWizardBuys(c)
	c.JSON(code, result)
}

func CkeyCharactersGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyCharacters(c)
	c.JSON(code, result)
}

func CharacterCkeysGET(c *gin.Context) {
	code, result := ckey_statistics.GetCharacterCkeys(c)
	c.JSON(code, result)
}

func CkeyRolesGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyRoles(c)
	c.JSON(code, result)
}
func AchievementsCkeysGET(c *gin.Context) {
	code, result := ckey_statistics.GetAchievementsCkey(c)
	c.JSON(code, result)
}

func AllRolesRoundsGET(c *gin.Context) {
	code, result := ckey_statistics.GetAllRolesRounds(c)
	c.JSON(code, result)
}