package ckey_statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

func GetCkeyUplinkBuys(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
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

	return 200, uplinkBuys
}
