package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
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
	var MapStatistics []*domain.MapStats
	r.Database.
		Preload("MapAttributes").
		Find(&MapStatistics)

	type simpleMapStats struct {
		MapName    string
		ServerID   string
		Attributes map[string]string
	}

	var maps []*simpleMapStats
	for _, stats := range MapStatistics {
		simpleMapStat := &simpleMapStats{
			MapName:    stats.MapName,
			ServerID:   stats.ServerID,
			Attributes: make(map[string]string),
		}
		for _, attribute := range stats.MapAttributes {
			simpleMapStat.Attributes[attribute.Name] = attribute.Value
		}
		maps = append(maps, simpleMapStat)
	}

	c.JSON(200, maps)
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
	_, processRoots, _, _, _ := getRoots(query, c)

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
		Preload("Factions", r.PreloadSelect("ID", "RootID", "Victory", "FactionName"), func(tx *gorm.DB) *gorm.DB {
			return tx.Where("faction_name <> ?", "Custom Squad")
		}).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "Victory", "RoleName"), func(tx *gorm.DB) *gorm.DB {
			return tx.Where("id in (select role_id from uplink_infos where id in (select uplink_info_id from uplink_purchases))")
		}).
		Preload("Factions.Members.UplinkInfo", r.PreloadSelect("ID", "RoleID")).
		Preload("Factions.Members.UplinkInfo.UplinkPurchases", r.PreloadSelect("UplinkInfoID", "ItemType", "Bundlename", "Cost"))
	_, processRoots, _, _, _ := getRoots(query, c)

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
