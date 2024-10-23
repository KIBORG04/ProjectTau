package ckey_statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/service/stats/ckey_statistics/crawler"
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
		Raw(`
		SELECT rolename, power_name, count(power_name) AS count,
		sum(case when faction_name in ? then f_victory
		when rolename in ? then r_victory
		else 1 end) AS wins
		FROM (
		select distinct on (f.root_id, roles.role_name, u.bundlename, f.faction_name) roles.role_name AS rolename, u.bundlename as power_name, f.faction_name as faction_name, roles.victory as r_victory, f.victory as f_victory
		from factions f
		join roles on roles.owner_id = f.id
		join uplink_infos i on roles.id = i.role_id
		join uplink_purchases u on i.id = u.uplink_info_id
		WHERE mind_ckey = ?
		) as t
		GROUP by rolename, power_name 
		ORDER BY count desc`, stats.TeamlRoles, stats.SoloRoles, player.Ckey).Scan(&uplinkBuys)

	if uplinkBuys == nil {
		return 404, map[string]string{
			"code":  "404",
			"error": "nothing found",
		}
	}

	for _, buy := range uplinkBuys {
		buy.Winrate = int(float32(buy.Wins) * 100 / float32(buy.Count))
	}

	return 200, uplinkBuys
}

func GetCkeyChangelingBuys(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	var changelingBuys []*AntagonistBuy

	r.Database.Raw(`
		SELECT rolename, power_name, count(power_name) AS count,
		sum(case when faction_name in ? then f_victory
		when rolename in ? then r_victory
		else 1 end) AS wins
		FROM (
		select distinct on (f.root_id, roles.role_name, u.power_name, f.faction_name) roles.role_name AS rolename, u.power_name as power_name, f.faction_name as faction_name, roles.victory as r_victory, f.victory as f_victory
		from factions f
		join roles on roles.owner_id = f.id
		join changeling_infos i on roles.id = i.role_id
		join changeling_purchases u on i.id = u.changeling_info_id
		WHERE mind_ckey = ?
		) as t
		GROUP by rolename, power_name 
		ORDER BY count desc`, stats.TeamlRoles, stats.SoloRoles, player.Ckey).Scan(&changelingBuys)

	if changelingBuys == nil {
		return 404, map[string]string{
			"code":  "404",
			"error": "nothing found",
		}
	}

	for _, buy := range changelingBuys {
		buy.Winrate = int(float32(buy.Wins) * 100 / float32(buy.Count))
	}

	return 200, changelingBuys
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

	r.Database.Raw(`
		SELECT rolename, power_name, count(power_name) AS count,
		sum(case when faction_name in ? then f_victory
		when rolename in ? then r_victory
		else 1 end) AS wins
		FROM (
		select distinct on (f.root_id, roles.role_name, u.power_name, f.faction_name) roles.role_name AS rolename, u.power_name as power_name, f.faction_name as faction_name, roles.victory as r_victory, f.victory as f_victory
		from factions f
		join roles on roles.owner_id = f.id
		join wizard_infos i on roles.id = i.role_id
		join wizard_purchases u on i.id = u.wizard_info_id
		WHERE mind_ckey = ?
		) as t
		GROUP by rolename, power_name 
		ORDER BY count desc`, stats.TeamlRoles, stats.SoloRoles, player.Ckey).Scan(&wizardBuys)

	if wizardBuys == nil {
		return 404, map[string]string{
			"code":  "404",
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
		return 404, map[string]string{
			"code":  "404",
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
		return 404, map[string]string{
			"code":  "404",
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

	var FoundChars []*struct {
		Name       string
		Similarity float32
	}

	r.Database.Raw(`
	select distinct mind_name as name, round(similarity(lower(mind_name), ?)::numeric * 100, 2) as similarity
	from roles
	where similarity(lower(mind_name), ?) > 0.3
	order by similarity desc;`, character.Name, character.Name).First(&FoundChars)

	if len(FoundChars) == 0 {
		return 404, map[string]string{
			"code":  "404",
			"error": "nothing found",
		}
	}

	return 200, FoundChars
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
		return 404, map[string]string{
			"code":  "404",
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
		return 404, map[string]string{
			"code":  "404",
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
		Date    string
		RoundId int
	}

	r.Database.Raw(`
		select root_id as round_id, name, title, a.desc, r.date
	    from achievements a
	    join roots r on r.round_id = a.root_id  
		where regexp_replace(lower(key), '['';()\\"&8^:$#№@_\s%]', '', 'g') = ?;
	`, player.Ckey).Scan(&achievementsInfo)

	if achievementsInfo == nil {
		return 404, map[string]string{
			"code":  "404",
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

func GetPlayerWithCrawler(c *gin.Context) (int, any) {
	player, err := stats.GetValidatePlayer(c)
	if err != nil {
		return 400, map[string]string{
			"code":  "400",
			"error": fmt.Sprint(err),
		}
	}

	playerStats := crawler.FetchPlayerStats(player.Ckey)

	if playerStats == nil {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	if len(playerStats.CrawlerStats) == 0 {
		return 400, map[string]string{
			"code":  "400",
			"error": "nothing found",
		}
	}

	type crawlerMinimized struct {
		ServerName string
		Minutes    uint
	}
	var crawlerStatsMinimized []crawlerMinimized

	for _, stat := range playerStats.CrawlerStats {
		crawlerStatsMinimized = append(crawlerStatsMinimized, crawlerMinimized{
			ServerName: stat.ServerName,
			Minutes:    stat.Minutes,
		})
	}

	return 200, crawlerStatsMinimized
}
