package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"sort"
	"ssstatistics/internal/config"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/keymap"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/utils"
	"strings"
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

type ServerCheckbox struct {
	Checkboxes []string `form:"server[]"`
}

func RootGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("Deaths", PreloadSelect("RootID", "Name", "SpecialRole", "AssignedRole")).
		Preload("CommunicationLogs", PreloadSelect("RootID", "Title", "Content", "Author")).
		Preload("Achievements", PreloadSelect("RootID", "Title", "Desc", "Key", "Name"))
	roots, processRoots, alphaRoots, betaRoots, gammaRoots := getRootsByCheckboxes(query, c)

	crewDeathsCount := make(keymap.MyMap[string, uint], 0)
	crewDeathsSum := 0

	roleDeathsCount := make(keymap.MyMap[string, uint], 0)
	roleDeathsSum := 0

	modesCount := make(keymap.MyMap[string, uint], 0)
	modesSum := 0

	var allAchievement []domain.Achievement

	var lastRoot *domain.Root
	for _, root := range processRoots {
		modesCount = keymap.AddElem(modesCount, root.Mode, 1)
		modesSum++
		for _, death := range root.Deaths {
			if isStationPlayer(death.AssignedRole, death.Name) {
				crewDeathsCount = keymap.AddElem(crewDeathsCount, death.AssignedRole, 1)
				crewDeathsSum++
			}
			if death.SpecialRole != "" {
				roleDeathsCount = keymap.AddElem(roleDeathsCount, death.SpecialRole, 1)
				roleDeathsSum++
			}
		}
		if lastRoot == nil || root.RoundID > lastRoot.RoundID {
			lastRoot = root
		}
		if len(root.Achievements) > 0 {
			allAchievement = append(allAchievement, root.Achievements...)
		}
	}

	sort.Stable(sort.Reverse(modesCount))
	sort.Stable(sort.Reverse(crewDeathsCount))
	sort.Stable(sort.Reverse(roleDeathsCount))

	var randRoot *domain.Root
	if len(processRoots) > 0 {
		randRoot = utils.Pick(processRoots)
	}

	var randComm domain.CommunicationLogs
	if randRoot != nil && len(randRoot.CommunicationLogs) > 0 {
		randComm = utils.Pick(randRoot.CommunicationLogs)
	}

	var lastAchievement domain.Achievement
	if len(allAchievement) > 0 {
		lastAchievement = utils.Pick(allAchievement)
	}

	var notNilLastRoot domain.Root
	if lastRoot != nil {
		notNilLastRoot = *lastRoot
	}

	return 200, "index.html", gin.H{
		"totalRounds": len(roots),
		"version":     notNilLastRoot.Version,
		"lastRound":   notNilLastRoot.RoundID,
		"lastDate":    notNilLastRoot.Date,
		"firstDate":   CurrentStatistics,

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
		Preload("LeaveStats", PreloadSelect("RootID", "Name", "AssignedRole", "LeaveTime", "LeaveType")).
		Preload("Factions", PreloadSelect("ID", "RootID", "FactionName", "Victory")).
		Preload("Factions.FactionObjectives", PreloadSelect("OwnerID", "Completed")).
		Preload("Factions.Members", PreloadSelect("ID", "OwnerID", "RoleName", "Victory")).
		Preload("Factions.Members.RoleObjectives", PreloadSelect("OwnerID", "Completed"))
	_, processRoots, _, _, _ := getRootsByCheckboxes(query, c)

	factionsSum := 0
	factionsCount := make(InfoSlice, 0)

	rolesSum := 0
	rolesCount := make(InfoSlice, 0)

	for _, root := range processRoots {
		var leavers uint
		for _, leaveStat := range root.LeaveStats {
			if isStationPlayer(leaveStat.AssignedRole, leaveStat.Name) && isRoundStartLeaver(leaveStat) {
				leavers++
			}
		}
		for _, faction := range root.Factions {
			foundInfo, ok := factionsCount.hasName(faction.FactionName)
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
				foundInfo, ok := rolesCount.hasName(role.RoleName)
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
		Preload("Factions", PreloadSelect("ID", "RootID", "Victory", "FactionName")).
		Preload("Factions.Members", PreloadSelect("ID", "OwnerID", "Victory", "RoleName")).
		Preload("Factions.Members.UplinkInfo", PreloadSelect("ID", "RoleID")).
		Preload("Factions.Members.UplinkInfo.UplinkPurchases", PreloadSelect("UplinkInfoID", "ItemType", "Bundlename", "Cost"))
	_, processRoots, _, _, _ := getRootsByCheckboxes(query, c)

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

			foundInfo, ok := infos.hasName(itemType)
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
				foundInfo, ok := uplinkRoles.hasName(roleName)
				if !ok {
					newUplinkInfo := &UplinkRoleInfo{
						Name:  roleName,
						Id:    ckey(roleName),
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
		Preload("Factions", PreloadSelect("ID", "RootID", "FactionName")).
		Preload("Factions.FactionObjectives", PreloadSelect("OwnerID", "Type", "Completed")).
		Preload("Factions.Members", PreloadSelect("ID", "OwnerID", "RoleName")).
		Preload("Factions.Members.RoleObjectives", PreloadSelect("OwnerID", "Type", "Completed"))
	_, processRoots, _, _, _ := getRootsByCheckboxes(query, c)

	objectiveHolders := make(InfoSlice, 0)

	addObjectiveInfo := func(owner *OwnerByObjectivesInfo, objective domain.Objectives) {
		owner.Count++
		var isWin uint
		if objective.Completed == ObjectiveWIN {
			isWin = 1
		}

		foundInfo, ok := owner.ObjectiveInfos.hasName(objective.Type)
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
				foundInfo, ok := objectiveHolders.hasName(faction.FactionName)
				var holderInfo *OwnerByObjectivesInfo
				if ok {
					holderInfo = (*foundInfo).(*OwnerByObjectivesInfo)
				} else {
					holderInfo = &OwnerByObjectivesInfo{
						Owner: faction.FactionName,
						Id:    ckey(faction.FactionName),
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
					foundInfo, ok := objectiveHolders.hasName(role.RoleName)
					var holderInfo *OwnerByObjectivesInfo
					if ok {
						holderInfo = (*foundInfo).(*OwnerByObjectivesInfo)
					} else {
						holderInfo = &OwnerByObjectivesInfo{
							Owner: role.RoleName,
							Id:    ckey(role.RoleName),
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
	_, processRoots, _, _, _ := getRootsByCheckboxes(query, c)

	return 200, "rounds.html", gin.H{
		"roots": processRoots,
	}
}

type PlayerTopInfo struct {
	Name  string
	Value uint
}

func (p PlayerTopInfo) GetName() string {
	return p.Name
}
func (p PlayerTopInfo) GetCount() uint {
	return p.Value
}

type TopInfo struct {
	Id          string
	Title       string
	NameColumn  string
	CountColumn string
	PlayersInfo InfoSlice

	ValuePostfix string
}

func (t *TopInfo) AddPlayerCount(name string) {
	t.ChangePlayerAndValue(name, 1, func(a uint, b uint) uint {
		return a + b
	})
}

func (t *TopInfo) SetPlayerAndValue(name string, value uint) {
	t.ChangePlayerAndValue(name, value, func(a uint, b uint) uint {
		return b
	})
}

func (t *TopInfo) ChangePlayerAndValue(name string, value uint, setRule func(uint, uint) uint) {
	foundInfo, ok := t.PlayersInfo.hasName(name)
	var holderInfo *PlayerTopInfo
	if ok {
		holderInfo = (*foundInfo).(*PlayerTopInfo)
		holderInfo.Value = setRule(holderInfo.Value, value)
	} else {
		holderInfo = &PlayerTopInfo{
			Name:  name,
			Value: value,
		}
		t.PlayersInfo = append(t.PlayersInfo, holderInfo)
	}
}

func TopsGET(c *gin.Context) (int, string, gin.H) {
	query := r.Database.
		Preload("Deaths", PreloadSelect("RootID", "MindName")).
		Preload("ManifestEntries", PreloadSelect("RootID", "AssignedRole", "Name")).
		Preload("LeaveStats", PreloadSelect("RootID", "AssignedRole", "Name", "LeaveTime", "LeaveType")).
		Preload("Score", PreloadSelect("ID", "RootID", "Richestkey", "Richestcash", "Dmgestkey", "Dmgestdamage")).
		Preload("Factions", PreloadSelect("RootID", "ID", "FactionName", "Victory")).
		Preload("Factions.Members", PreloadSelect("ID", "OwnerID", "MindCkey", "MindName", "RoleName", "Victory"))

	_, processRoots, _, _, _ := getRootsByCheckboxes(query, c)

	tops := make([]*TopInfo, 0)
	gamemodeTops := make([]*TopInfo, 0)

	hasId := func(slice []*TopInfo, id string) bool {
		for _, topInfo := range slice {
			if topInfo.Id == id {
				return true
			}
		}
		return false
	}

	getTopById := func(slice []*TopInfo, id string) *TopInfo {
		for _, topInfo := range slice {
			if topInfo.Id == id {
				return topInfo
			}
		}
		return nil
	}

	deathTop := &TopInfo{
		Id:          "deaths",
		Title:       "Смертей",
		NameColumn:  "Имя",
		CountColumn: "Количество",
	}
	tops = append(tops, deathTop)

	roundsPlayedTop := &TopInfo{
		Id:          "zadrots",
		Title:       "Задротов",
		NameColumn:  "Имя",
		CountColumn: "Раундов",
	}
	tops = append(tops, roundsPlayedTop)

	leaversTop := &TopInfo{
		Id:          "leavers",
		Title:       "Ливеров",
		NameColumn:  "Имя",
		CountColumn: "Количество",
	}
	tops = append(tops, leaversTop)

	richestTop := &TopInfo{
		Id:           "rich",
		Title:        "Богатейших",
		NameColumn:   "Имя",
		CountColumn:  "Денег",
		ValuePostfix: "$",
	}
	tops = append(tops, richestTop)

	damagedTop := &TopInfo{
		Id:          "damaged",
		Title:       "Избитых",
		NameColumn:  "Имя",
		CountColumn: "Урон",
	}
	tops = append(tops, damagedTop)

	ckeyAntagsPlayed := make(map[string]map[string]map[string]uint)
	for _, root := range processRoots {
		for _, death := range root.Deaths {
			// wtf
			if death.MindName == "Unknown" || death.MindName == "unknown" {
				continue
			}
			deathTop.AddPlayerCount(death.MindName)
		}
		for _, entry := range root.ManifestEntries {
			if isStationPlayer(entry.AssignedRole, entry.Name) {
				roundsPlayedTop.AddPlayerCount(entry.Name)
			}
		}
		for _, stat := range root.LeaveStats {
			if isStationPlayer(stat.AssignedRole, stat.Name) && isRoundStartLeaver(stat) {
				leaversTop.AddPlayerCount(stat.Name)
			}
		}

		for _, faction := range root.Factions {
			if slices.Contains(teamlRoles, faction.FactionName) && !hasId(gamemodeTops, ckey(faction.FactionName)) {
				title := faction.FactionName
				if value, ok := shortModeName[faction.FactionName]; ok {
					title = value
				}
				gamemodeTop := &TopInfo{
					Id:           ckey(faction.FactionName),
					Title:        title,
					NameColumn:   "Ckey",
					CountColumn:  "Winrate",
					ValuePostfix: "%",
				}
				gamemodeTops = append(gamemodeTops, gamemodeTop)
			}

			processedCkeys := make([]string, 0, len(faction.Members))
			for _, role := range faction.Members {
				if slices.Contains(processedCkeys, role.MindCkey) {
					continue
				}

				// uniq mode checks
				if faction.FactionName == "Cult Of Blood" &&
					(utils.IsMobName.FindString(role.MindName) != "" || strings.Contains(role.MindName, "familiar")) {
					continue
				}

				if slices.Contains(soloRoles, role.RoleName) && !hasId(gamemodeTops, ckey(role.RoleName)) {
					title := role.RoleName
					if value, ok := shortModeName[role.RoleName]; ok {
						title = value
					}
					gamemodeTop := &TopInfo{
						Id:           ckey(role.RoleName),
						Title:        title,
						NameColumn:   "Ckey",
						CountColumn:  "Winrate",
						ValuePostfix: "%",
					}
					gamemodeTops = append(gamemodeTops, gamemodeTop)
				}

				if _, ok := ckeyAntagsPlayed[role.MindCkey]; !ok {
					antagMap := make(map[string]map[string]uint)
					ckeyAntagsPlayed[role.MindCkey] = antagMap
				}

				if slices.Contains(teamlRoles, faction.FactionName) {
					if _, ok := ckeyAntagsPlayed[role.MindCkey][faction.FactionName]; !ok {
						infoMap := map[string]uint{
							"Victory": uint(faction.Victory),
							"Count":   1,
						}
						ckeyAntagsPlayed[role.MindCkey][faction.FactionName] = infoMap
					} else {
						ckeyAntagsPlayed[role.MindCkey][faction.FactionName]["Victory"] += uint(faction.Victory)
						ckeyAntagsPlayed[role.MindCkey][faction.FactionName]["Count"] += 1
					}
				}

				if slices.Contains(soloRoles, role.RoleName) {
					roleWin := uint(role.Victory)
					if faction.FactionName == "Shadowlings" {
						roleWin = uint(faction.Victory)
					}
					if _, ok := ckeyAntagsPlayed[role.MindCkey][role.RoleName]; !ok {
						infoMap := map[string]uint{
							"Victory": roleWin,
							"Count":   1,
						}
						ckeyAntagsPlayed[role.MindCkey][role.RoleName] = infoMap
					} else {
						ckeyAntagsPlayed[role.MindCkey][role.RoleName]["Victory"] += roleWin
						ckeyAntagsPlayed[role.MindCkey][role.RoleName]["Count"] += 1
					}
				}

				processedCkeys = append(processedCkeys, role.MindCkey)
			}
		}

		richestTop.SetPlayerAndValue(root.Score.Richestkey, uint(root.Score.Richestcash))
		damagedTop.SetPlayerAndValue(root.Score.Dmgestkey, uint(root.Score.Dmgestdamage))
	}

	for player, antagInfo := range ckeyAntagsPlayed {
		for antag, antagOptions := range antagInfo {
			if antagOptions["Count"] > 10 {
				antagTop := getTopById(gamemodeTops, ckey(antag))
				antagTop.SetPlayerAndValue(player, uint(float32(antagOptions["Victory"]*100)/float32(antagOptions["Count"])))
			}
		}
	}
	// remove useless positions
	for _, top := range tops {
		if len(top.PlayersInfo) == 0 {
			tops = utils.RemoveElem(tops, top)
		}
		sort.Sort(sort.Reverse(top.PlayersInfo))
		if len(top.PlayersInfo) > 10 {
			top.PlayersInfo = slices.Delete(top.PlayersInfo, 10, len(top.PlayersInfo))
		}
	}
	for _, top := range gamemodeTops {
		if len(top.PlayersInfo) == 0 {
			gamemodeTops = utils.RemoveElem(gamemodeTops, top)
		}
		sort.Sort(sort.Reverse(top.PlayersInfo))
		if len(top.PlayersInfo) > 10 {
			top.PlayersInfo = slices.Delete(top.PlayersInfo, 10, len(top.PlayersInfo))
		}
	}

	return 200, "tops.html", gin.H{
		"topSlice":         tops,
		"gamemodeTopSlice": gamemodeTops,
	}
}
