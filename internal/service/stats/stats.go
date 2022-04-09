package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"sort"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/keymap"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/utils"
	"strings"
)

type ServerCheckbox struct {
	Checkboxes []string `form:"server[]"`
}

func RootPOST(c *gin.Context) {
	setCheckboxStates(c)
	RootGET(c)
}

func RootGET(c *gin.Context) {
	checkboxStates := getCheckboxStates(c)
	roots, processRoots, alphaRoots, betaRoots, gammaRoots := getRootsByCheckboxes([]string{"Deaths"}, checkboxStates)

	crewDeathsCount := make(keymap.MyMap[string, uint], 0)
	crewDeathsSum := 0

	roleDeathsCount := make(keymap.MyMap[string, uint], 0)
	roleDeathsSum := 0

	modesCount := make(keymap.MyMap[string, uint], 0)
	modesSum := 0

	var lastRoot *domain.Root
	for _, root := range processRoots {
		modesCount = keymap.AddElem(modesCount, root.Mode, 1)
		modesSum++
		for _, death := range root.Deaths {
			if slices.Contains(stationPositions, death.AssignedRole) && utils.IsDrone.FindString(death.Name) == "" {
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
	}

	sort.Stable(sort.Reverse(modesCount))
	sort.Stable(sort.Reverse(crewDeathsCount))
	sort.Stable(sort.Reverse(roleDeathsCount))

	var links domain.Link
	str := fmt.Sprintf("%%%d%%", lastRoot.RoundID)
	r.Database.Where("link LIKE ?", str).First(&links)

	c.HTML(200, "index.html", gin.H{
		"totalRounds": len(roots),
		"version":     lastRoot.Version,
		"lastRound":   lastRoot.RoundID,
		"lastDate":    links.Date,

		"alphaRounds": len(alphaRoots),
		"betaRounds":  len(betaRoots),
		"gammaRounds": len(gammaRoots),

		"modesCount":      modesCount,
		"modesSum":        modesSum,
		"crewDeathsCount": crewDeathsCount,
		"crewDeathsSum":   crewDeathsSum,
		"roleDeathsCount": roleDeathsCount,
		"roleDeathsSum":   roleDeathsSum,

		"serverCheckboxes": checkboxStates,
	})
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

func GamemodesPOST(c *gin.Context) {
	setCheckboxStates(c)
	GamemodesGET(c)
}

func GamemodesGET(c *gin.Context) {
	checkboxStates := getCheckboxStates(c)
	_, processRoots, _, _, _ := getRootsByCheckboxes([]string{"Factions.FactionObjectives", "Factions.Members.RoleObjectives"}, checkboxStates)

	factionsSum := 0
	factionsCount := make(InfoSlice, 0)

	rolesSum := 0
	rolesCount := make(InfoSlice, 0)

	for _, root := range processRoots {
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

	c.HTML(200, "gamemodes.html", gin.H{
		"factionsSum":      factionsSum,
		"factionsCount":    factionsCount,
		"rolesSum":         rolesSum,
		"rolesCount":       rolesCount,
		"serverCheckboxes": checkboxStates,
	})
}

func Cult(c *gin.Context) {
	render := charts.Cult()

	c.HTML(200, "chart.html", gin.H{
		"charts": render,
	})
}

type UplinkInfo struct {
	Name      string
	Count     uint
	Type      string
	Wins      uint
	Winrate   uint
	TotalCost uint
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

func UplinkPOST(c *gin.Context) {
	setCheckboxStates(c)
	UplinkGET(c)
}

func UplinkGET(c *gin.Context) {
	checkboxStates := getCheckboxStates(c)
	_, processRoots, _, _, _ := getRootsByCheckboxes([]string{"Factions.Members.UplinkInfo.UplinkPurchases"}, checkboxStates)

	uplinkRoles := make(InfoSlice, 0)

	addUplinkInfo := func(infos InfoSlice, counter *uint, purchase *domain.UplinkPurchases, faction *domain.Factions, role *domain.Role) InfoSlice {
		foundInfo, ok := infos.hasName(purchase.ItemType)
		var isWin uint
		if faction != nil {
			isWin = uint(faction.Victory)
		} else {
			isWin = uint(role.Victory)
		}
		if !ok {
			infos = append(infos, &UplinkInfo{
				Name:      purchase.Bundlename,
				Count:     1,
				Type:      purchase.ItemType,
				Wins:      isWin,
				Winrate:   isWin * 100,
				TotalCost: uint(purchase.Cost),
			})
		} else {
			uplinkInfo := (*foundInfo).(*UplinkInfo)
			uplinkInfo.Count++
			uplinkInfo.Wins += isWin
			uplinkInfo.Winrate = uplinkInfo.Wins * 100 / uplinkInfo.Count
			uplinkInfo.TotalCost += uint(purchase.Cost)

		}
		*counter++
		return infos
	}

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			for _, role := range faction.Members {
				for _, purchase := range role.UplinkInfo.UplinkPurchases {
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
							Id:    strings.ReplaceAll(strings.ToLower(roleName), " ", ""),
							Count: 1,
						}
						s := StatInfo(newUplinkInfo)
						foundInfo = &s
						uplinkRoles = append(uplinkRoles, newUplinkInfo)
					} else {
						(*foundInfo).(*UplinkRoleInfo).Count++
					}
					newUplinkInfo := (*foundInfo).(*UplinkRoleInfo)
					newUplinkInfo.UplinkInfos = addUplinkInfo(newUplinkInfo.UplinkInfos, &newUplinkInfo.Count, &purchase, useFaction, &role)
				}
			}
		}
	}

	for _, role := range uplinkRoles {
		uplinkRoleInfo := role.(*UplinkRoleInfo)
		sort.Sort(sort.Reverse(uplinkRoleInfo.UplinkInfos))
	}

	c.HTML(200, "uplink.html", gin.H{
		"uplinkPurchases":  uplinkRoles,
		"serverCheckboxes": checkboxStates,
	})
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

func ObjectivesPOST(c *gin.Context) {
	setCheckboxStates(c)
	ObjectivesGET(c)
}

func ObjectivesGET(c *gin.Context) {
	checkboxStates := getCheckboxStates(c)
	_, processRoots, _, _, _ := getRootsByCheckboxes([]string{"Factions.FactionObjectives", "Factions.Members.RoleObjectives"}, checkboxStates)

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
						Id:    strings.ReplaceAll(strings.ToLower(faction.FactionName), " ", ""),
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
							Id:    strings.ReplaceAll(strings.ToLower(role.RoleName), " ", ""),
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

	c.HTML(200, "objectives.html", gin.H{
		"objectiveHolders": objectiveHolders,
		"serverCheckboxes": checkboxStates,
	})
}
