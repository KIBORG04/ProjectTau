package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm/clause"
	"image"
	"sort"
	"ssstatistics/internal/config"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/keymap"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/utils"
)

func BasicGET(c *gin.Context, f func(*gin.Context) (int, string, gin.H)) {
	baseContext := gin.H{
		"baseUrl":          config.Config.BaseUrl,
		"serverCheckboxes": getCheckboxStates(c),
	}

	code, file, context := f(c)

	for k, v := range baseContext {
		context[k] = v
	}

	c.HTML(code, file, context)
}

func BasicPOST(c *gin.Context, f func(*gin.Context) (int, string, gin.H)) {
	setCheckboxStates(c)
	BasicGET(c, f)
}

// RootGET TODO: remove keymap.KeyMap
func RootGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("Deaths", r.PreloadSelect("RootID", "Name", "SpecialRole", "AssignedRole", "DeathX", "DeathY", "DeathZ")).
		Preload("Explosions", r.PreloadSelect("RootID", "EpicenterX", "EpicenterY", "EpicenterZ")).
		Preload("CommunicationLogs", r.PreloadSelect("RootID", "Title", "Content", "Author")).
		Preload("Achievements", r.PreloadSelect("RootID", "Title", "Desc", "Key", "Name"))
	roots, processRoots, alphaRoots, betaRoots, gammaRoots := getRoots(query, c)

	crewDeathsCount := make(keymap.MyMap[string, uint], 0)
	crewDeathsSum := 0

	roleDeathsCount := make(keymap.MyMap[string, uint], 0)
	roleDeathsSum := 0

	modesCount := make(keymap.MyMap[string, uint], 0)
	modesSum := 0

	deathCoords := make([]image.Point, 0, crewDeathsSum)
	explosionCoords := make([]image.Point, 0)

	var allAchievement []domain.Achievement
	var allAnnounces []domain.CommunicationLogs

	var lastRoot *domain.Root
	for _, root := range processRoots {
		modesCount = keymap.AddElem(modesCount, root.Mode, 1)
		modesSum++
		for _, death := range root.Deaths {
			if IsStationPlayer(death.AssignedRole, death.Name) {
				crewDeathsCount = keymap.AddElem(crewDeathsCount, death.AssignedRole, 1)
				crewDeathsSum++
			}
			if death.SpecialRole != "" {
				roleDeathsCount = keymap.AddElem(roleDeathsCount, death.SpecialRole, 1)
				roleDeathsSum++
			}
			if death.DeathZ == 2 && root.Map == "Box Station" {
				deathCoords = append(deathCoords, image.Point{int(death.DeathX), int(death.DeathY)})
			}
		}
		if lastRoot == nil || root.RoundID > lastRoot.RoundID {
			lastRoot = root
		}
		if len(root.Achievements) > 0 {
			allAchievement = append(allAchievement, root.Achievements...)
		}
		for _, log := range root.CommunicationLogs {
			allAnnounces = append(allAnnounces, log)
		}
		for _, explosion := range root.Explosions {
			if explosion.EpicenterZ == 2 {
				explosionCoords = append(explosionCoords, image.Point{int(explosion.EpicenterX), int(explosion.EpicenterY)})
			}
		}
	}

	sort.Stable(sort.Reverse(modesCount))
	sort.Stable(sort.Reverse(crewDeathsCount))
	sort.Stable(sort.Reverse(roleDeathsCount))

	var randComm domain.CommunicationLogs
	if len(allAnnounces) > 0 {
		randComm = utils.Pick(allAnnounces)
	}

	var lastAchievement domain.Achievement
	if len(allAchievement) > 0 {
		lastAchievement = utils.Pick(allAchievement)
	}

	var notNilLastRoot domain.Root
	if lastRoot != nil {
		notNilLastRoot = *lastRoot
	}

	// TODO
	// heatmap.Create("deaths", deathCoords)
	// heatmap.Create("explosions", explosionCoords)

	return 200, "index.html", gin.H{
		"totalRounds": len(roots),
		"version":     notNilLastRoot.Version,
		"lastRound":   notNilLastRoot.RoundID,
		"lastDate":    utils.TrimPGDate(notNilLastRoot.Date),
		"firstDate":   CurrentStatisticsDate,

		"alphaRounds": len(alphaRoots),
		"betaRounds":  len(betaRoots),
		"gammaRounds": len(gammaRoots),

		"modesCount":      modesCount,
		"modesSum":        modesSum,
		"crewDeathsCount": crewDeathsCount,
		"crewDeathsSum":   crewDeathsSum,
		"roleDeathsCount": roleDeathsCount,
		"roleDeathsSum":   roleDeathsSum,

		"randComm":        randComm,
		"lastAchievement": lastAchievement,
	}
}

type FactionInfo struct {
	Name                string
	Count               uint
	Wins                uint
	Winrate             uint
	Members             uint
	TotalObjectives     uint
	CompletedObjectives uint
	PercentObjectives   uint
	Leavers             uint
	AvgLeavers          float32
}

func (b FactionInfo) GetName() string {
	return b.Name
}
func (b FactionInfo) GetCount() uint {
	return b.Count
}

type RolesInfo struct {
	Name                string
	Count               uint
	Wins                uint
	Winrate             uint
	TotalObjectives     uint
	CompletedObjectives uint
	PercentObjectives   uint
}

func (b RolesInfo) GetName() string {
	return b.Name
}
func (b RolesInfo) GetCount() uint {
	return b.Count
}

// functions declaration in another function doesn't support generics
func completedObjectives[T any](objectives []T) uint {
	var completed uint
	for _, objective := range objectives {
		switch t := any(objective).(type) {
		case domain.RoleObjectives:
			if t.Completed == ObjectiveWIN {
				completed++
			}
		case domain.FactionObjectives:
			if t.Completed == ObjectiveWIN {
				completed++
			}
		}
	}
	return completed
}

func GamemodesGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("LeaveStats", r.PreloadSelect("RootID", "Name", "AssignedRole", "LeaveTime", "LeaveType")).
		Preload("Factions", r.PreloadSelect("ID", "RootID", "FactionName", "Victory")).
		Preload("Factions.FactionObjectives", r.PreloadSelect("OwnerID", "Completed")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "RoleName", "Victory")).
		Preload("Factions.Members.RoleObjectives", r.PreloadSelect("OwnerID", "Completed"))
	_, processRoots, _, _, _ := getRoots(query, c)

	factionsSum := 0
	factionsCount := make(InfoSlice, 0)

	rolesSum := 0
	rolesCount := make(InfoSlice, 0)

	for _, root := range processRoots {
		var leavers uint
		for _, leaveStat := range root.LeaveStats {
			if IsStationPlayer(leaveStat.AssignedRole, leaveStat.Name) && IsRoundStartLeaver(leaveStat) {
				leavers++
			}
		}
		for _, faction := range root.Factions {
			foundInfo, ok := factionsCount.HasName(faction.FactionName)
			if !ok {
				factionsCount = append(factionsCount, &FactionInfo{
					Name:                faction.FactionName,
					Count:               1,
					Wins:                uint(faction.Victory),
					Members:             uint(len(faction.Members)),
					Winrate:             uint(faction.Victory * 100),
					TotalObjectives:     uint(len(faction.FactionObjectives)),
					CompletedObjectives: completedObjectives(faction.FactionObjectives),
					PercentObjectives:   completedObjectives(faction.FactionObjectives) * 100 / utils.Max(uint(len(faction.FactionObjectives)), 1),
					Leavers:             leavers,
					AvgLeavers:          float32(leavers),
				})
			} else {
				factionInfo := (*foundInfo).(*FactionInfo)
				factionInfo.Count++
				factionInfo.Members += uint(len(faction.Members))
				factionInfo.Wins += uint(faction.Victory)
				factionInfo.Winrate = factionInfo.Wins * 100 / factionInfo.Count
				factionInfo.TotalObjectives += uint(len(faction.FactionObjectives))
				factionInfo.CompletedObjectives += completedObjectives(faction.FactionObjectives)
				factionInfo.PercentObjectives = factionInfo.CompletedObjectives * 100 / utils.Max(factionInfo.TotalObjectives, 1)
				factionInfo.Leavers += leavers
				factionInfo.AvgLeavers = float32(factionInfo.Leavers) / float32(factionInfo.Count)
			}
			factionsSum++

			for _, role := range faction.Members {
				foundInfo, ok := rolesCount.HasName(role.RoleName)
				if !ok {
					rolesCount = append(rolesCount, &RolesInfo{
						Name:                role.RoleName,
						Count:               1,
						Wins:                uint(role.Victory),
						Winrate:             uint(role.Victory) * 100,
						TotalObjectives:     uint(len(role.RoleObjectives)),
						CompletedObjectives: completedObjectives(role.RoleObjectives),
						PercentObjectives:   completedObjectives(role.RoleObjectives) * 100 / utils.Max(uint(len(role.RoleObjectives)), 1),
					})
				} else {
					roleInfo := (*foundInfo).(*RolesInfo)
					roleInfo.Count++
					roleInfo.Wins += uint(role.Victory)
					roleInfo.Winrate = roleInfo.Wins * 100 / roleInfo.Count
					roleInfo.TotalObjectives += uint(len(role.RoleObjectives))
					roleInfo.CompletedObjectives += completedObjectives(role.RoleObjectives)
					roleInfo.PercentObjectives = roleInfo.CompletedObjectives * 100 / utils.Max(roleInfo.TotalObjectives, 1)
				}
				rolesSum++
			}

		}
	}

	sort.Sort(sort.Reverse(factionsCount))
	sort.Sort(sort.Reverse(rolesCount))

	return 200, "gamemodes.html", gin.H{
		"factionsSum":   factionsSum,
		"factionsCount": factionsCount,
		"rolesSum":      rolesSum,
		"rolesCount":    rolesCount,
	}
}

func Cult(c *gin.Context) (int, string, gin.H) {
	render := charts.Cult()

	return 200, "chart.html", gin.H{
		"charts": render,
	}
}

type UplinkInfo struct {
	Name       string
	Count      uint
	TotalCount uint
	Type       string
	Wins       uint
	Winrate    uint
	TotalCost  uint
}

func (b UplinkInfo) GetName() string {
	return b.Type
}
func (b UplinkInfo) GetCount() uint {
	return b.Count
}

type UplinkRoleInfo struct {
	Name        string
	Count       uint
	Id          string
	UplinkInfos InfoSlice
}

func (b UplinkRoleInfo) GetName() string {
	return b.Name
}
func (b UplinkRoleInfo) GetCount() uint {
	return b.Count
}

func UplinkGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("Factions", r.PreloadSelect("ID", "RootID", "Victory", "FactionName")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "Victory", "RoleName")).
		Preload("Factions.Members.UplinkInfo", r.PreloadSelect("ID", "RoleID")).
		Preload("Factions.Members.UplinkInfo.UplinkPurchases", r.PreloadSelect("UplinkInfoID", "ItemType", "Bundlename", "Cost"))
	_, processRoots, _, _, _ := getRoots(query, c)

	uplinkRoles := make(InfoSlice, 0)

	addUplinkInfo := func(infos InfoSlice, purchases []domain.UplinkPurchases, faction *domain.Factions, role *domain.Role) InfoSlice {
		var isWin uint
		if faction != nil {
			isWin = uint(faction.Victory)
		} else {
			isWin = uint(role.Victory)
		}

		processed := make([]string, 0, len(purchases))

		for _, purchase := range purchases {
			itemType := purchase.ItemType
			if itemType == "" {
				itemType = purchase.Bundlename
			}
			itemName := purchase.Bundlename
			// бандлу с рандомным лутом ставится такое же название, что и виду коробки
			// покупка "рандомного итема" имеет тот же тайп, но цену в 0
			if itemType == "/obj/item/weapon/storage/box/syndicate" {
				if purchase.Cost > 0 {
					itemType = itemName
				} else {
					itemName = "Random Item"
					itemType = "Random Item"
				}
			}

			foundInfo, ok := infos.HasName(itemType)
			if !ok {
				infos = append(infos, &UplinkInfo{
					Name:       itemName,
					Count:      1,
					TotalCount: 1,
					Type:       itemType,
					Wins:       isWin,
					Winrate:    isWin * 100,
					TotalCost:  uint(purchase.Cost),
				})
			} else {
				uplinkInfo := (*foundInfo).(*UplinkInfo)
				uplinkInfo.TotalCount++
				uplinkInfo.TotalCost += uint(purchase.Cost)
				if !slices.Contains(processed, itemType) {
					uplinkInfo.Count++
					uplinkInfo.Wins += isWin
					uplinkInfo.Winrate = uplinkInfo.Wins * 100 / uplinkInfo.Count
				}
			}
			processed = append(processed, itemType)
		}
		return infos
	}

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			for _, role := range faction.Members {
				if len(role.UplinkInfo.UplinkPurchases) == 0 {
					continue
				}
				roleName := role.RoleName

				var useFaction *domain.Factions
				if faction.FactionName == "Syndicate Operatives" || faction.FactionName == "Revolution" {
					roleName = faction.FactionName
					useFaction = &faction
				}
				foundInfo, ok := uplinkRoles.HasName(roleName)
				if !ok {
					newUplinkInfo := &UplinkRoleInfo{
						Name:  roleName,
						Id:    Ckey(roleName),
						Count: 1,
					}
					s := StatInfo(newUplinkInfo)
					foundInfo = &s
					uplinkRoles = append(uplinkRoles, newUplinkInfo)
				} else {
					(*foundInfo).(*UplinkRoleInfo).Count++
				}

				newUplinkInfo := (*foundInfo).(*UplinkRoleInfo)
				newUplinkInfo.UplinkInfos = addUplinkInfo(newUplinkInfo.UplinkInfos, role.UplinkInfo.UplinkPurchases, useFaction, &role)
			}
		}
	}

	return 200, "uplink.html", gin.H{
		"uplinkPurchases": uplinkRoles,
	}
}

type ObjectiveInfo struct {
	Type    string
	Count   uint
	Wins    uint
	Winrate uint
}

func (b ObjectiveInfo) GetName() string {
	return b.Type
}
func (b ObjectiveInfo) GetCount() uint {
	return b.Count
}

type OwnerByObjectivesInfo struct {
	Owner          string
	Count          uint
	Id             string
	ObjectiveInfos InfoSlice
}

func (b OwnerByObjectivesInfo) GetName() string {
	return b.Owner
}
func (b OwnerByObjectivesInfo) GetCount() uint {
	return b.Count
}

func ObjectivesGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("Factions", r.PreloadSelect("ID", "RootID", "FactionName")).
		Preload("Factions.FactionObjectives", r.PreloadSelect("OwnerID", "Type", "Completed")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "RoleName")).
		Preload("Factions.Members.RoleObjectives", r.PreloadSelect("OwnerID", "Type", "Completed"))
	_, processRoots, _, _, _ := getRoots(query, c)

	objectiveHolders := make(InfoSlice, 0)

	addObjectiveInfo := func(owner *OwnerByObjectivesInfo, objective domain.Objectives) {
		owner.Count++
		var isWin uint
		if objective.Completed == ObjectiveWIN {
			isWin = 1
		}

		foundInfo, ok := owner.ObjectiveInfos.HasName(objective.Type)
		var objectiveInfo *ObjectiveInfo
		if ok {
			objectiveInfo = (*foundInfo).(*ObjectiveInfo)
			objectiveInfo.Count++
			objectiveInfo.Wins += isWin
			objectiveInfo.Winrate = objectiveInfo.Wins * 100 / objectiveInfo.Count
		} else {
			objectiveInfo = &ObjectiveInfo{
				Type:    objective.Type,
				Count:   1,
				Wins:    isWin,
				Winrate: isWin * 100,
			}
			owner.ObjectiveInfos = append(owner.ObjectiveInfos, objectiveInfo)
		}
	}

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			if len(faction.FactionObjectives) > 0 {
				foundInfo, ok := objectiveHolders.HasName(faction.FactionName)
				var holderInfo *OwnerByObjectivesInfo
				if ok {
					holderInfo = (*foundInfo).(*OwnerByObjectivesInfo)
				} else {
					holderInfo = &OwnerByObjectivesInfo{
						Owner: faction.FactionName,
						Id:    Ckey(faction.FactionName),
					}
					objectiveHolders = append(objectiveHolders, holderInfo)
				}
				for _, objective := range faction.FactionObjectives {
					// IDK HOW
					if objective.Type == "" {
						continue
					}
					addObjectiveInfo(holderInfo, domain.Objectives(objective))
				}
			}

			for _, role := range faction.Members {
				if len(role.RoleObjectives) > 0 {
					foundInfo, ok := objectiveHolders.HasName(role.RoleName)
					var holderInfo *OwnerByObjectivesInfo
					if ok {
						holderInfo = (*foundInfo).(*OwnerByObjectivesInfo)
					} else {
						holderInfo = &OwnerByObjectivesInfo{
							Owner: role.RoleName,
							Id:    Ckey(role.RoleName),
						}
						objectiveHolders = append(objectiveHolders, holderInfo)
					}
					for _, objective := range role.RoleObjectives {
						addObjectiveInfo(holderInfo, domain.Objectives(objective))
					}
				}
			}
		}
	}

	return 200, "objectives.html", gin.H{
		"objectiveHolders": objectiveHolders,
	}
}

func RoundGET(c *gin.Context) (int, string, gin.H) {
	roundId := c.Param("id")
	if utils.RoundId.FindString(roundId) == "" {
		return 400, "error.html", gin.H{
			"error": fmt.Sprintf("%s is not round id.", roundId),
		}
	}
	root, err := r.EagerFindByRoundId(roundId)
	if err != nil {
		return 404, "error.html", gin.H{
			"error": fmt.Sprintf("%s round not found.", roundId),
		}
	}
	return 200, "round.html", gin.H{
		"root": root,
	}
}

func RoundsGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.Order("round_id DESC").Limit(100)
	_, processRoots, _, _, _ := getRoots(query, c)

	return 200, "rounds.html", gin.H{
		"roots": processRoots,
	}
}

func TopsGET(c *gin.Context) (int, string, gin.H) {
	var tops []*domain.Top

	r.Database.
		Preload(clause.Associations).
		Order("Title").
		Find(&tops)

	return 200, "tops.html", gin.H{
		"topSlice": tops,
	}
}

func MmrGET(c *gin.Context) (int, string, gin.H) {
	return 200, "mmr.html", gin.H{}
}

func MapsGET(c *gin.Context) (int, string, gin.H) {
	return 200, "maps.html", gin.H{}
}

func FeedbackGET(c *gin.Context) (int, string, gin.H) {
	return 200, "feedback.html", gin.H{}
}
