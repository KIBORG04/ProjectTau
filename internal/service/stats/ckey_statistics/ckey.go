package ckey_statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/utils"
	"strings"
)

// TODO: возможно надо перенести валидацию приходящих значений в контроллер-часть. Тут же должны тупо выполняться запросы в БД

type AntagonistBuy struct {
	Rolename  string
	PowerName string
	Wins      int
	Winrate   int
	Count     int
}

func GetCkeyUplinkBuys(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var uplinkBuys []*AntagonistBuy

	r.Database.
		Select(`  roles.role_name 									   AS rolename,
						u.bundlename 										   AS power_name, 
					    count(u.bundlename) 								   AS count,
						sum(case when factions.faction_name in ? then factions.victory
						when roles.role_name in ? then roles.victory
						else 1 end) 										   AS wins`, stats.TeamlRoles, stats.SoloRoles).
		Table("factions").
		Joins("join roles on roles.owner_id = factions.id").
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

	for _, buy := range uplinkBuys {
		buy.Winrate = int(float32(buy.Wins) * 100 / float32(buy.Count))
	}

	return 200, uplinkBuys
}

func GetCkeyChanglingBuys(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var changlingBuys []*AntagonistBuy

	r.Database.
		Select(`  roles.role_name 									   AS rolename,
						u.power_name, 
					    count(u.power_name) 								   AS count,
						sum(case when factions.faction_name in ? then factions.victory
						when roles.role_name in ? then roles.victory
						else 1 end) 										   AS wins`, stats.TeamlRoles, stats.SoloRoles).
		Table("factions").
		Joins("join roles on roles.owner_id = factions.id").
		Joins("join changeling_infos i on roles.id = i.role_id").
		Joins("join changeling_purchases u on i.id = u.changeling_info_id").
		Where("mind_ckey = ?", player.Ckey).
		Group("roles.role_name, u.power_name").
		Order("count desc").
		Find(&changlingBuys)

	if changlingBuys == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	for _, buy := range changlingBuys {
		buy.Winrate = int(float32(buy.Wins) * 100 / float32(buy.Count))
	}

	return 200, changlingBuys
}

func GetCkeyWizardBuys(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var wizardBuys []*AntagonistBuy

	r.Database.
		Select(`  roles.role_name 									   AS rolename,
						u.power_name, 
					    count(u.power_name) 								   AS count,
						sum(case when factions.faction_name in ? then factions.victory
						when roles.role_name in ? then roles.victory
						else 1 end) 										   AS wins`, stats.TeamlRoles, stats.SoloRoles).
		Table("factions").
		Joins("join roles on roles.owner_id = factions.id").
		Joins("join wizard_infos i on roles.id = i.role_id").
		Joins("join wizard_purchases u on i.id = u.wizard_info_id").
		Where("mind_ckey = ?", player.Ckey).
		Group("roles.role_name, u.power_name").
		Order("count desc").
		Find(&wizardBuys)

	if wizardBuys == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	for _, buy := range wizardBuys {
		buy.Winrate = int(float32(buy.Wins) * 100 / float32(buy.Count))
	}

	return 200, wizardBuys
}

func FindSimilaryCkey(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var FoundCkey struct {
		FoundCkey string
	}

	r.Database.Raw(`
	select mind_ckey as found_ckey, similarity(mind_ckey, ?) as sim
	from roles
	where similarity(mind_ckey, ?) > 0.4
	order by sim desc;`, player.Ckey, player.Ckey).First(&FoundCkey)

	if FoundCkey.FoundCkey == "" {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, FoundCkey
}

func GetCkeyCharacters(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var characters []*struct {
		MindName string
		Count    int
	}

	r.Database.Raw(`
	select mind_name, count(1) as count
	from roles
	where mind_ckey = ?
	  and mind_name not like '%(%)'
	  and mind_name not like 'homunculus%'
	  and mind_name not like 'Syndicate Robot-%'
	  and role_name not in ('Abductor Scientist', 'Abductor Agent', 'Cortical Borer')
	  and mind_name <> ''
	  and mind_name <> 'unknown'
	group by mind_name
	order by count desc;`, player.Ckey).Scan(&characters)

	if characters == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, characters
}

func FindSimilaryCharacter(c *gin.Context) (int, any) {
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
	character.Name = strings.ToLower(character.Name)

	var FoundChar struct {
		Name string
	}

	r.Database.Raw(`
	select mind_name as name, similarity(lower(mind_name), ?) as sim
	from roles
	where similarity(lower(mind_name), ?) > 0.3
	order by sim desc;`, character.Name, character.Name).First(&FoundChar)

	if FoundChar.Name == "" {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, FoundChar
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
	character.Name = strings.ToLower(character.Name)

	var ckeys []*struct {
		MindCkey string
		Count    int
	}

	r.Database.Raw(`
	select mind_ckey, count(1) as count
	from roles
	where lower(mind_name) =  ?
	group by mind_ckey
	order by count desc;
	`, character.Name).Scan(&ckeys)

	if ckeys == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, ckeys
}

func GetCkeyRoles(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var rolesInfo []*struct {
		RoleName string
		Count    int
		Wins     int
	}

	r.Database.Raw(`
		select role_name,
			   count(1) as count,
			   sum(case when faction_name in ? then faction_victory
						when role_name in ? then role_victory
						else 1
				   end) as wins
		from (select distinct factions.root_id,
							  factions.faction_name,
							  factions.victory as faction_victory,
							  r.role_name,
							  r.victory as role_victory
			  from factions
			  join roles r on factions.id = r.owner_id
			  where r.mind_ckey = ?) as t
		group by role_name
		order by role_name;
		`, stats.TeamlRoles, stats.SoloRoles, player.Ckey).Scan(&rolesInfo)

	if rolesInfo == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, rolesInfo
}

func GetAchievementsCkey(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var achievementsInfo []*struct {
		Name    string
		Title   string
		Desc    string
		RoundId int
	}

	r.Database.Raw(`
		select root_id as round_id, name, title, "desc" 
	    from achievements
		where regexp_replace(lower(key), '['';()\\"&8^:$#№@_\s%]', '', 'g') = ?;
	`, player.Ckey).Scan(&achievementsInfo)

	if achievementsInfo == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, achievementsInfo
}

func GetAllRolesRounds(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var allAntagsInfo []*struct {
		RoundId     int
		Date        string
		FactionName string
		RoleName    string
		Win         int
	}

	r.Database.Raw(`
	select round_id, 
	       date, 
	       f.faction_name, 
	       r.role_name, 
	       (case when f.faction_name in ? then f.victory
                 when r.role_name in ? then r.victory
                 else 1
    		end) as win
	from roots
	join factions f on f.root_id = round_id
	join roles r on r.owner_id = f.id
	where r.mind_ckey = ?
	group by f.faction_name, r.role_name, round_id, win
	order by round_id desc;

	`, stats.TeamlRoles, stats.SoloRoles, player.Ckey).Scan(&allAntagsInfo)

	if allAntagsInfo == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	for _, info := range allAntagsInfo {
		info.Date = utils.TrimPGDate(info.Date)
	}

	return 200, allAntagsInfo
}

func GetCkeyMMR(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var MMR []*struct {
		Mmr int
	}

	r.Database.Raw(`
	select mmr
	from players
	where ckey = ?;`, player.Ckey).Scan(&MMR)

	if MMR == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	return 200, MMR
}
