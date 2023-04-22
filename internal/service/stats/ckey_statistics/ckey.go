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

	if uplinkBuys == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, uplinkBuys
}

func GetCkeyCharacters(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var characters []struct {
		MindName string
		Count    int
	}

	r.Database.Raw(`
	select mind_name, count(1)
	from roles
	where mind_ckey = ?
	  and mind_name not like '%(%)'
	  and mind_name not like 'homunculus%'
	  and role_name not in ('Abductor Scientist', 'Abductor Agent', 'Cortical Borer')
	group by mind_name;`, player.Ckey).Scan(&characters)

	if characters == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, characters
}

func GetCharacterCkeys(c *gin.Context) (int, any) {
	type Character struct {
		Name string `form:"name"`
	}
	var character Character
	err := c.BindQuery(&character)
	if err != nil {
		return 400, err
	}
	if character.Name == "" {
		return 400, map[string]string{
			"code":  "400",
			"error": "name not entered in query",
		}
	}

	var ckeys []struct {
		MindCkey string
		Count    int
	}

	r.Database.Raw(`
	select mind_ckey, count(1) as count
	from roles
	where mind_name = ?
	group by mind_ckey
	order by count desc
	`, character.Name).Scan(&ckeys)

	if ckeys == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, ckeys
}
