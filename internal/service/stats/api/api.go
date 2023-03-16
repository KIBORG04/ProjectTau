package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"ssstatistics/internal/bots/telegram"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/service/stats/last_phrase"
	"strings"
	"time"
)

func MmrGET(c *gin.Context) {
	type mmr struct {
		Ckey string
		MMR  int32
	}

	var mmrs []*mmr
	player, _ := getValidatePlayer(c)
	if player != (*Player)(nil) {
		r.Database.
			Model(&domain.Player{}).
			Select("Ckey", "MMR").
			Where("ckey = ?", player.Ckey).
			Find(&mmrs)
	} else {
		r.Database.
			Model(&domain.Player{}).
			Select("Ckey", "MMR").
			Find(&mmrs)
	}

	c.JSON(200, mmrs)
}

func MapsGET(c *gin.Context) {
	var mapStatistics []*struct {
		Name           string
		Server         string
		Count          int
		Crewscore      float32
		Stuffshipped   float32
		Stuffharvested float32
		Oremined       float32
		Researchdone   float32
		Powerloss      float32
		Mess           float32
		Meals          float32
		Nuked          float32
		Recantags      float32
		Crewescaped    float32
		Crewdead       float32
		Crewtotal      float32
		Crewsurvived   float32
		Foodeaten      float32
		Clownabuse     float32
		Duration       float32
	}

	query := r.Database.Table("roots r").
		Select(`
		case r.map 
			when 'Falcon Station (Snowy)' then 'Falcon Station'
			when 'Gamma Station (Snowy)' then 'Gamma Station'
			when 'Prometheus Station (Snowy)' then 'Prometheus Station'
			else r.map
		end as name,
		r.server_address as server,
		avg(s.crewscore) as crewscore,
		avg(s.stuffshipped) as stuffshipped,
		avg(s.stuffharvested) as stuffharvested,
		avg(s.oremined) as oremined,
		avg(s.researchdone) as researchdone,
		avg(s.powerloss) as powerloss,
		avg(s.mess) as mess,
		avg(s.meals) as meals,
		avg(s.nuked) as nuked,
		avg(s.rec_antags) as RecAntags,
		avg(s.crew_escaped) as crewescaped,
		avg(s.crew_dead) as crewdead,
		avg(s.crew_total) as crewtotal,
		avg(s.crew_survived) as crewsurvived,
		avg(s.foodeaten) as foodeaten,
		avg(s.clownabuse) as clownabuse,
		count(r.map) as count,
		avg(case when duration = '' then 3600
           else split_part(duration, ':', 1)::int * 3600 + split_part(duration, ':', 2)::int * 60
        end) as duration
		`).
		Joins("join scores s on s.root_id = r.round_id").
		Group("name, r.server_address")
	stats.ApplyDBQueryByDate(query, c)
	query.Find(&mapStatistics)

	c.JSON(200, mapStatistics)
}

func SendFeedback(c *gin.Context) {
	type FeedbackForm struct {
		Username string `json:"username"`
		Text     string `json:"text"`
	}

	var form FeedbackForm
	if err := c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, "Некорректный запрос.")
		return
	}

	if len(form.Username) < 3 || len(form.Username) > 50 {
		c.JSON(http.StatusBadRequest, "В имени должно быть от 3х до 50 символов.")
		return
	}

	if len(form.Text) < 10 || len(form.Text) > 255 {
		c.JSON(http.StatusBadRequest, "В сообщение должно быть от 10 до 255 символов..")
		return
	}

	msg := fmt.Sprintf(`
	*%s*
	%s`, form.Username, form.Text)

	err := telegram.Bot.Send(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Сообщение отправлено.")
}

func HeatmapsGET(c *gin.Context) {
	type Request struct {
		// explosions, deaths
		Type          string `form:"type"`
		MapResolution string `form:"mapresolution"`
		MapName       string `form:"mapname"`
	}

	type Response struct {
		MapBase64 string
		Error     string
	}

	var query Request
	err := c.BindQuery(&query)
	if err != nil {
		c.JSON(200, Response{Error: "The request is incorrect"})
		return
	}

	c.JSON(200, Response{
		MapBase64: "SoSi He-x bibu",
		Error:     "228 Server shlet tebya нахуй."})
}

func ChanglingGET(c *gin.Context) {
	var (
		RolesCount []struct {
			RoleName string
			Count    int
		}

		ChanglingStats []struct {
			RoleName  string
			PowerName string
			PowerType string
			Count     int
			Wins      int
			TotalCost int
			Winrate   float32
		}
	)

	var (
		startDate, endDate, _ = stats.GetValidDates(c)
		s1, s2, s3            = stats.GetChosenServers(c)
	)

	r.Database.Raw(`
	select r.role_name, count(1)
	from roles r
	join changeling_infos ci on ci.role_id = r.id
	join factions f on f.id = r.owner_id 
	join roots root on root.round_id = f.root_id
	where (root.date between ? and ?)
	and server_address = ? OR server_address = ? OR server_address = ?
	group by r.role_name;
	`,
		startDate, endDate,
		s1, s2, s3).Scan(&RolesCount)

	r.Database.Raw(`
	select  role_name,
	        power_name,
			split_part(power_type, '/', -1) as power_type, 
			count(1) as count,
			sum(victory) as wins, 
			sum(cost) as total_cost,
			sum(victory)*100 / count(1)::real as winrate
	from (
		select r.victory, r.role_name, cp.power_type, cp.power_name, cp.cost
		from roles r
		join changeling_infos ci on ci.role_id = r.id
		join changeling_purchases cp on cp.changeling_info_id = ci.id
		join factions f on f.id = r.owner_id 
		join roots root on root.round_id = f.root_id
		where (root.date between ? and ?)
		and server_address = ? OR server_address = ? OR server_address = ?
		) as a
		group by role_name, power_type, power_name;
	`,
		startDate, endDate,
		s1, s2, s3).Scan(&ChanglingStats)

	type AbilityInfo struct {
		Name      string
		Count     int
		Wins      int
		Winrate   int
		TotalCost int
	}

	type Info struct {
		Count              int
		ChanglingAbilities map[string]*AbilityInfo
	}

	roleAbilitiesMap := make(map[string]*Info)

	for _, role := range RolesCount {
		roleAbilitiesMap[role.RoleName] = &Info{
			Count:              role.Count,
			ChanglingAbilities: make(map[string]*AbilityInfo),
		}
	}

	for _, stat := range ChanglingStats {
		roleAbilitiesMap[stat.RoleName].ChanglingAbilities[stat.PowerType] = &AbilityInfo{
			Name:      stat.PowerName,
			Count:     stat.Count,
			Wins:      stat.Wins,
			Winrate:   int(stat.Winrate),
			TotalCost: stat.TotalCost,
		}
	}

	c.JSON(200, roleAbilitiesMap)
}

func UplinkGET(c *gin.Context) {
	type (
		UplinkInfo struct {
			Name       string
			Count      uint
			TotalCount uint
			Wins       uint
			Winrate    uint
			TotalCost  uint
		}

		UplinkRoleInfo struct {
			Name        string
			Count       uint
			UplinkInfos map[string]*UplinkInfo
		}
	)

	query := r.Database.
		Preload("Factions",
			r.PreloadSelect("ID", "RootID", "Victory", "FactionName"),
			func(tx *gorm.DB) *gorm.DB {
				return tx.Where("faction_name <> ?", "Custom Squad")
			}).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "Victory", "RoleName"),
			func(tx *gorm.DB) *gorm.DB {
				return tx.Where("id in (select role_id from uplink_infos where id in (select uplink_info_id from uplink_purchases))")
			}).
		Preload("Factions.Members.UplinkInfo", r.PreloadSelect("ID", "RoleID")).
		Preload("Factions.Members.UplinkInfo.UplinkPurchases", r.PreloadSelect("UplinkInfoID", "ItemType", "Bundlename", "Cost"))
	processRoots := stats.GetRoots(query, c)

	uplinkRolesMap := make(map[string]*UplinkRoleInfo, 0)

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			if faction.FactionName == "Custom Squad" {
				continue
			}
			for _, role := range faction.Members {
				if len(role.UplinkInfo.UplinkPurchases) == 0 {
					continue
				}
				roleName := role.RoleName

				var useFaction bool
				if faction.FactionName == "Syndicate Operatives" || faction.FactionName == "Revolution" {
					roleName = faction.FactionName
					useFaction = true
				}
				uplinkRole, ok := uplinkRolesMap[stats.Ckey(roleName)]
				if !ok {
					uplinkRole = &UplinkRoleInfo{
						Name:        roleName,
						Count:       1,
						UplinkInfos: make(map[string]*UplinkInfo),
					}
					uplinkRolesMap[stats.Ckey(roleName)] = uplinkRole
				} else {
					uplinkRole.Count++
				}
				var isWin uint
				if useFaction {
					isWin = uint(faction.Victory)
				} else {
					isWin = uint(role.Victory)
				}

				processed := make([]string, 0, len(role.UplinkInfo.UplinkPurchases))
				for _, purchase := range role.UplinkInfo.UplinkPurchases {
					itemType := purchase.ItemType
					itemName := purchase.Bundlename
					if itemType == "" {
						itemType = stats.Ckey(purchase.Bundlename)
					} else if itemType != "/obj/item/weapon/storage/box/syndicate" {
						splitType := strings.Split(purchase.ItemType, "/")
						itemType = splitType[len(splitType)-2] + "_" + splitType[len(splitType)-1]
					} else if itemType == "/obj/item/weapon/storage/box/syndicate" {
						// бандлу с рандомным лутом ставится такое же название, что и виду коробки
						// покупка "рандомного итема" имеет тот же тайп, но цену в 0
						if purchase.Cost > 0 {
							itemType = stats.Ckey(itemName)
						} else {
							itemName = "Random Item"
							itemType = "RandomItem"
						}
					}

					uplinkPurchase, ok := uplinkRole.UplinkInfos[itemType]
					if !ok {
						uplinkPurchase = &UplinkInfo{
							Name:       itemName,
							Count:      1,
							TotalCount: 1,
							Wins:       isWin,
							Winrate:    isWin * 100,
							TotalCost:  uint(purchase.Cost),
						}
						uplinkRole.UplinkInfos[itemType] = uplinkPurchase
					} else {
						uplinkPurchase.TotalCount++
						uplinkPurchase.TotalCost += uint(purchase.Cost)
						if !slices.Contains(processed, itemType) {
							uplinkPurchase.Count++
							uplinkPurchase.Wins += isWin
							uplinkPurchase.Winrate = uplinkPurchase.Wins * 100 / uplinkPurchase.Count
						}
					}
					processed = append(processed, itemType)
				}
			}
		}
	}

	c.JSON(200, uplinkRolesMap)
}

func RandomAnnounceGET(c *gin.Context) {
	var randComm domain.CommunicationLogs
	r.Database.Model(&randComm).
		Select("Title", "Content", "Author").
		Where("type not like 'fax%'").
		Order("random()").
		Limit(1).
		Find(&randComm)

	var announce struct {
		Title   string
		Content string
		Author  string
	}

	announce.Title = randComm.Title
	announce.Content = template.HTMLEscapeString(randComm.Content)
	announce.Author = randComm.Author

	c.JSON(200, announce)

}

func RandomAchievementGET(c *gin.Context) {
	var randAchievement domain.Achievement
	r.Database.Model(&randAchievement).Select("Title", "Desc", "Key", "Name").Order("random()").Limit(1).
		Find(&randAchievement)

	var achievement struct {
		Title string
		Desc  string
		Key   string
		Name  string
	}

	achievement.Title = randAchievement.Title
	achievement.Desc = template.HTMLEscapeString(randAchievement.Desc)
	achievement.Key = randAchievement.Key
	achievement.Name = randAchievement.Name

	c.JSON(200, achievement)

}

func RandomLastPhraseGET(c *gin.Context) {
	c.JSON(200, last_phrase.GetRandomLastPhrase())
}

func RandomFlavorGET(c *gin.Context) {
	var randEntries domain.ManifestEntries
	r.Database.Model(&randEntries).
		Select("Name", "Species", "Gender", "Age", "Flavor").
		Where("flavor <> '' AND char_length(flavor) > 10 AND name NOT LIKE 'syndicate drone%'").
		Order("random()").
		Limit(1).
		Find(&randEntries)

	var lastPhrase struct {
		Name    string
		Species string
		Gender  string
		Age     uint
		Flavor  string
	}

	lastPhrase.Name = randEntries.Name
	lastPhrase.Species = randEntries.Species
	lastPhrase.Gender = randEntries.Gender
	lastPhrase.Age = randEntries.Age
	lastPhrase.Flavor = template.HTMLEscapeString(randEntries.Flavor)

	c.JSON(200, lastPhrase)
}

func ModeWinratesByMonthGET(c *gin.Context) {
	type FactionOrRole struct {
		FactionName string `form:"faction"`
		RoleName    string `form:"role"`
	}

	var query FactionOrRole
	if err := c.BindQuery(&query); err != nil || (query.FactionName == "" && query.RoleName == "") {
		c.JSON(http.StatusBadRequest, "Некорректный запрос.")
		return
	}

	if query.FactionName != "" && query.RoleName != "" {
		c.JSON(http.StatusBadRequest, "Либо фракция, либо роль")
		return
	}

	type SqlParams struct {
		Antag    string
		DateFrom string
		DateTo   string
	}

	output := make(map[string]int)

	if query.FactionName != "" {
		var allFactions []string
		r.Database.Model(&domain.Factions{}).Distinct("FactionName").Find(&allFactions)

		if !slices.Contains(allFactions, query.FactionName) {
			c.JSON(http.StatusBadRequest, "Фракция не найдена.")
			return
		}

		startDate, _ := time.Parse("2006-01-02", stats.ModeStartDate)
		endDate := time.Now()

		params := SqlParams{
			Antag: query.FactionName,
		}
		for (startDate.Month() != endDate.Month()+1) || (startDate.Year() != endDate.Year()) {
			dateFromString := startDate.Format("2006-01-02")
			startDate = startDate.AddDate(0, 1, 0)

			params.DateFrom = dateFromString
			params.DateTo = startDate.Format("2006-01-02")

			var winrate float32
			r.Database.Raw(`
		select SUM(victory)::real * 100 / COUNT(id)::real as winrate
		from factions
		join roots r on r.round_id = factions.root_id
		where faction_name = @Antag and date >= @DateFrom and date <= @DateTo;
		`, params).Scan(&winrate)

			output[dateFromString] = int(winrate)
		}

	} else {
		var allRoles []string
		r.Database.Model(&domain.Role{}).Distinct("RoleName").Find(&allRoles)

		if !slices.Contains(allRoles, query.RoleName) {
			c.JSON(http.StatusBadRequest, "Фракция не найдена.")
			return
		}

		startDate, _ := time.Parse("2006-01-02", stats.ModeStartDate)
		endDate := time.Now()

		params := SqlParams{
			Antag: query.RoleName,
		}
		for (startDate.Month() != endDate.Month()+1) || (startDate.Year() != endDate.Year()) {
			dateFromString := startDate.Format("2006-01-02")
			startDate = startDate.AddDate(0, 1, 0)

			params.DateFrom = dateFromString
			params.DateTo = startDate.Format("2006-01-02")

			var winrate float32
			r.Database.Raw(`
		select SUM(roles.victory)::real * 100 / COUNT(roles.id)::real as winrate
		from roles
		join factions f on f.id = roles.owner_id
		join roots r on r.round_id = f.root_id
		where role_name = @Antag and date >= @DateFrom and date <= @DateTo;
		`, params).Scan(&winrate)

			output[dateFromString] = int(winrate)
		}
	}

	c.JSON(http.StatusOK, output)
}
