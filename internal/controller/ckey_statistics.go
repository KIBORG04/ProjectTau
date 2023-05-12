package controller

import (
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/service/stats/ckey_statistics"
)

func FinderGET(c *gin.Context) (int, string, gin.H) {
	return 200, "finder.html", gin.H{}
}

func PlayerStatisticGET(c *gin.Context) (int, string, gin.H) {
	return 200, "player.html", gin.H{}
}

func NotAPlayerGET(c *gin.Context) (int, string, gin.H) {
	return 200, "error.html", gin.H{
		"error": "Введи игрока нормально, без абузов",
	}
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

func TryFindCkeyGET(c *gin.Context) {
	code, result := ckey_statistics.FindSimilaryCkey(c)
	c.JSON(code, result)
}

func CkeyCharactersGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyCharacters(c)
	c.JSON(code, result)
}

func TryFindCharacterGET(c *gin.Context) {
	code, result := ckey_statistics.FindSimilaryCharacter(c)
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

func CkeyMMRGET(c *gin.Context) {
	code, result := ckey_statistics.GetCkeyMMR(c)
	c.JSON(code, result)
}
