package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

type Player struct {
	Ckey string `form:"ckey"`
}

func getValidatePlayer(c *gin.Context) (*Player, error) {
	var player Player
	err := c.BindQuery(&player)
	if err != nil {
		return nil, err
	}
	if player.Ckey == "" {
		return nil, fmt.Errorf("ckey not entered in query")
	}
	player.Ckey = stats.Ckey(player.Ckey)
	return &player, nil
}

func CkeyUplinkBuysGET(c *gin.Context) {
	player, err := getValidatePlayer(c)
	if err != nil {
		c.JSON(400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		})
		return
	}

	var uplinkBuys []struct {
		Rolename   string
		Bundlename string
		Count      int
	}

	r.Database.
		Select("roles.role_name as rolename,u.bundlename, count(u.bundlename) as count").
		Table("roles").
		Joins("join uplink_infos i on roles.id = i.role_id").
		Joins("join uplink_purchases u on i.id = u.uplink_info_id").
		Where("mind_ckey = ?", player.Ckey).
		Group("roles.role_name, u.bundlename").
		Order("count desc").
		Find(&uplinkBuys)

	c.JSON(200, uplinkBuys)
}
