package stats

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
	"strings"
)

func ApiMmrGET(c *gin.Context) {
	var players []*domain.Player
	r.Database.
		Select("Ckey", "MMR").
		Find(&players)

	type mmr struct {
		Ckey string
		MMR  int32
	}

	var mmrs []*mmr
	for _, player := range players {
		mmrs = append(mmrs, &mmr{
			Ckey: player.Ckey,
			MMR:  player.MMR,
		})
	}

	c.JSON(200, mmrs)
}

func ApiMapsGET(c *gin.Context) {
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
		Select("r.map as name, " +
			"r.server_address as server, " +
			"avg(s.crewscore) as crewscore, " +
			"avg(s.stuffshipped) as stuffshipped, " +
			"avg(s.stuffharvested) as stuffharvested, " +
			"avg(s.oremined) as oremined, " +
			"avg(s.researchdone) as researchdone, " +
			"avg(s.powerloss) as powerloss, " +
			"avg(s.mess) as mess, " +
			"avg(s.meals) as meals, " +
			"avg(s.nuked) as nuked, " +
			"avg(s.rec_antags) as RecAntags, " +
			"avg(s.crew_escaped) as crewescaped, " +
			"avg(s.crew_dead) as crewdead, " +
			"avg(s.crew_total) as crewtotal, " +
			"avg(s.crew_survived) as crewsurvived, " +
			"avg(s.foodeaten) as foodeaten, " +
			"avg(s.clownabuse) as clownabuse, " +
			"count(r.map) as count").
		Joins("join scores s on s.root_id = r.round_id").
		Group("r.map, r.server_address")
	applyDBQueryByDate(query, c)
	query.Find(&mapStatistics)

	var roots []*domain.Root
	query = r.Database.Select("duration", "map", "server_address")
	applyDBQueryByDate(query, c)
	query.Find(&roots)

	for _, root := range roots {
		for _, statistic := range mapStatistics {
			if statistic.Name == root.Map && statistic.Server == root.ServerAddress {
				roundTime, err := ParseRoundTime(root.Duration)
				if err == nil {
					statistic.Duration += float32(roundTime.ToSeconds())
				} else {
					statistic.Duration += float32(3600)
				}
			}
		}
	}

	for _, statistic := range mapStatistics {
		statistic.Duration = statistic.Duration / float32(statistic.Count)
	}

	c.JSON(200, mapStatistics)
}

func ApiSendFeedback(c *gin.Context) {
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
	Name: *%s*
	
Text: %s`, form.Username, form.Text)

	err := telegram.Bot.Send(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Сообщение отправлено.")
}

func ApiHeatmapsGET(c *gin.Context) {
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

type (
	changlingRole struct {
		Count              uint
		ChanglingAbilities map[string]*changlingAbilities
	}

	changlingAbilities struct {
		Name      string
		Count     uint
		Wins      uint
		Winrate   uint
		TotalCost uint
	}
)

func ApiChanglingGET(c *gin.Context) {
	query := r.Database.
		Preload("Factions", r.PreloadSelect("ID", "RootID", "Victory", "FactionName")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "Victory", "RoleName")).
		Preload("Factions.Members.ChangelingInfo", r.PreloadSelect("ID", "RoleID")).
		Preload("Factions.Members.ChangelingInfo.ChangelingPurchase", r.PreloadSelect("ChangelingInfoID", "PowerType", "PowerName", "Cost"))
	processRoots := getRoots(query, c)

	roleAbilitiesMap := make(map[string]*changlingRole)

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			for _, member := range faction.Members {
				if len(member.ChangelingInfo.ChangelingPurchase) == 0 {
					continue
				}
				role, ok := roleAbilitiesMap[member.RoleName]
				if !ok {
					role = &changlingRole{
						Count:              1,
						ChanglingAbilities: make(map[string]*changlingAbilities),
					}
					roleAbilitiesMap[member.RoleName] = role
				} else {
					role.Count++
				}

				for _, purchase := range member.ChangelingInfo.ChangelingPurchase {
					powerType := strings.Split(purchase.PowerType, "/")
					abilityType := powerType[len(powerType)-1]
					ability, ok := role.ChanglingAbilities[abilityType]
					if !ok {
						role.ChanglingAbilities[abilityType] = &changlingAbilities{
							Name:      purchase.PowerName,
							Count:     1,
							Wins:      uint(member.Victory),
							Winrate:   uint(member.Victory * 100),
							TotalCost: uint(purchase.Cost),
						}
					} else {
						ability.Count++
						ability.Wins += uint(member.Victory)
						ability.Winrate = ability.Wins * 100 / ability.Count
						ability.TotalCost += uint(purchase.Cost)
					}
				}
			}
		}
	}

	c.JSON(200, roleAbilitiesMap)
}

type (
	UplinkRoleInfo struct {
		Name        string
		Count       uint
		UplinkInfos map[string]*UplinkInfo
	}

	UplinkInfo struct {
		Name       string
		Count      uint
		TotalCount uint
		Wins       uint
		Winrate    uint
		TotalCost  uint
	}
)

func ApiUplinkGET(c *gin.Context) {
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
	processRoots := getRoots(query, c)

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
				uplinkRole, ok := uplinkRolesMap[Ckey(roleName)]
				if !ok {
					uplinkRole = &UplinkRoleInfo{
						Name:        roleName,
						Count:       1,
						UplinkInfos: make(map[string]*UplinkInfo),
					}
					uplinkRolesMap[Ckey(roleName)] = uplinkRole
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
						itemType = Ckey(purchase.Bundlename)
					} else if itemType != "/obj/item/weapon/storage/box/syndicate" {
						splitType := strings.Split(purchase.ItemType, "/")
						itemType = splitType[len(splitType)-1]
					} else if itemType == "/obj/item/weapon/storage/box/syndicate" {
						// бандлу с рандомным лутом ставится такое же название, что и виду коробки
						// покупка "рандомного итема" имеет тот же тайп, но цену в 0
						if purchase.Cost > 0 {
							itemType = Ckey(itemName)
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

func ApiRandomAnnounceGET(c *gin.Context) {
	var randComm domain.CommunicationLogs
	r.Database.Model(&randComm).Select("Title", "Content", "Author").Order("random()").Limit(1).
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

func ApiRandomAchievementGET(c *gin.Context) {
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
