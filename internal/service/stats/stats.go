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

type FactionsInfos []*FactionsInfo

func (info FactionsInfos) Len() int {
	return len(info)
}

func (info FactionsInfos) Less(i, j int) bool {
	return info[i].Winrate < info[j].Winrate
}

func (info FactionsInfos) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}

func (info FactionsInfos) hasName(name string) (*FactionsInfo, bool) {
	for i := 0; i < len(info); i++ {
		if info[i].Name == name {
			return info[i], true
		}
	}
	return nil, false
}

type FactionsInfo struct {
	Name                string
	Count               uint
	Wins                uint
	Members             uint
	Winrate             uint
	TotalObjectives     uint
	CompletedObjectives uint
	PercentObjectives   uint
}

type RolesInfos []*RolesInfo

func (info RolesInfos) Len() int {
	return len(info)
}

func (info RolesInfos) Less(i, j int) bool {
	return info[i].Winrate < info[j].Winrate
}

func (info RolesInfos) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}

func (info RolesInfos) hasName(name string) (*RolesInfo, bool) {
	for i := 0; i < len(info); i++ {
		if info[i].Name == name {
			return info[i], true
		}
	}
	return nil, false
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

func Gamemodes(c *gin.Context) {
	checkboxStates := getCheckboxStates(c)
	_, processRoots, _, _, _ := getRootsByCheckboxes([]string{"Factions.FactionObjectives", "Factions.Members.RoleObjectives"}, checkboxStates)

	factionsSum := 0
	factionsCount := make(FactionsInfos, 0)

	rolesSum := 0
	rolesCount := make(RolesInfos, 0)

	for _, root := range processRoots {
		for _, faction := range root.Factions {
			foundInfo, ok := factionsCount.hasName(faction.FactionName)
			if !ok {
				factionsCount = append(factionsCount, &FactionsInfo{
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
				foundInfo.Count++
				foundInfo.Members += uint(len(faction.Members))
				foundInfo.Wins += uint(faction.Victory)
				foundInfo.Winrate = foundInfo.Wins * 100 / foundInfo.Count
				foundInfo.TotalObjectives += uint(len(faction.FactionObjectives))
				foundInfo.CompletedObjectives += completedObjectives(faction.FactionObjectives)
				foundInfo.PercentObjectives = foundInfo.CompletedObjectives * 100 / utils.Max(foundInfo.TotalObjectives, 1)
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
					foundInfo.Count++
					foundInfo.Wins += uint(role.Victory)
					foundInfo.Winrate = foundInfo.Wins * 100 / foundInfo.Count
					foundInfo.TotalObjectives += uint(len(role.RoleObjectives))
					foundInfo.CompletedObjectives += completedObjectives(role.RoleObjectives)
					foundInfo.PercentObjectives = foundInfo.CompletedObjectives * 100 / utils.Max(foundInfo.TotalObjectives, 1)
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
