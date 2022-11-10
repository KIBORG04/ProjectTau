package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
	"ssstatistics/internal/config"
	"ssstatistics/internal/domain"
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
	type Group struct {
		Name  string
		Count int
	}

	var modesCount []Group
	r.Database.Model(&domain.Root{}).
		Select("mode as name, count(mode) as count").
		Where("mode <> ''").
		Group("mode").
		Order("count desc").
		Scan(&modesCount)

	var modesSum int
	for _, group := range modesCount {
		modesSum += group.Count
	}

	var crewDeathsCount []Group
	r.Database.Model(&domain.Deaths{}).
		Select("assigned_role as name, count(assigned_role) as count").
		Where("assigned_role in (?)", stationPositions).
		Where("name not like 'maintenance drone%'").
		Group("assigned_role").
		Order("count desc").
		Scan(&crewDeathsCount)

	var crewDeathsSum int
	for _, group := range crewDeathsCount {
		crewDeathsSum += group.Count
	}

	var roleDeathsCount []Group
	r.Database.Model(&domain.Deaths{}).
		Select("special_role as name, count(special_role) as count").
		Where("special_role <> ''").
		Group("special_role").
		Order("count desc").
		Scan(&roleDeathsCount)

	var roleDeathsSum int
	for _, group := range roleDeathsCount {
		roleDeathsSum += group.Count
	}

	/* TODO:
	deathCoords := make([]image.Point, 0, crewDeathsSum)
	explosionCoords := make([]image.Point, 0)

	for _, root := range processRoots {
		for _, death := range root.Deaths {
			if death.DeathZ == 2 && root.Map == "Box Station" {
				deathCoords = append(deathCoords, image.Point{int(death.DeathX), int(death.DeathY)})
			}
		}
		for _, explosion := range root.Explosions {
			if explosion.EpicenterZ == 2 && root.Map == "Box Station" && explosion.DevastationRange == 0 && explosion.HeavyImpactRange == 0 && explosion.LightImpactRange == 2 && explosion.FlashRange == 3 {
				explosionCoords = append(explosionCoords, image.Point{int(explosion.EpicenterX), int(explosion.EpicenterY)})
			}
		}

	Preload("Explosions", r.PreloadSelect("RootID", "EpicenterX", "EpicenterY", "EpicenterZ", "DevastationRange", "HeavyImpactRange", "LightImpactRange", "FlashRange")).
	go heatmap.Create(c.Request.Context(), "deaths", deathCoords)
	go heatmap.Create(c.Request.Context(), "explosions", explosionCoords)
	}*/

	var randComm domain.CommunicationLogs
	r.Database.Model(&randComm).Select("Title", "Content", "Author").Order("random()").Limit(1).
		Find(&randComm)

	var lastAchievement domain.Achievement
	r.Database.Model(&lastAchievement).Select("Title", "Desc", "Key", "Name").Order("random()").Limit(1).
		Find(&lastAchievement)

	var lastRoot domain.Root
	r.Database.Model(&lastRoot).Select("Version", "RoundID", "Date").Order("round_id desc").Limit(1).
		Find(&lastRoot)

	var rootsStatistics []struct {
		Server string
		Count  int
	}

	r.Database.Model(&domain.Root{}).
		Select("server_address as server, count(server_address) as count").
		Where("server_address <> ''").
		Group("server_address").
		Scan(&rootsStatistics)

	var (
		totalRoots      int
		totalAlphaRoots int
		totalBetaaRoots int
		totalGammaRoots int
	)
	for _, root := range rootsStatistics {
		totalRoots += root.Count
		switch root.Server {
		case ServerAlphaAddress:
			totalAlphaRoots += root.Count
		case ServerBetaAddress:
			totalBetaaRoots += root.Count
		case ServerGammaAddress:
			totalGammaRoots += root.Count
		}
	}

	return 200, "index.html", gin.H{
		"totalRounds": totalRoots,
		"version":     lastRoot.Version,
		"lastRound":   lastRoot.RoundID,
		"lastDate":    utils.TrimPGDate(lastRoot.Date),
		"firstDate":   CurrentStatisticsDate,

		"alphaRounds": totalAlphaRoots,
		"betaRounds":  totalBetaaRoots,
		"gammaRounds": totalGammaRoots,

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
	processRoots := getRoots(query, c)

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

func UplinkGET(c *gin.Context) (int, string, gin.H) {
	return 200, "uplink.html", gin.H{}
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
		Preload("Factions.Members",
			r.PreloadSelect("ID", "OwnerID", "RoleName"),
			func(tx *gorm.DB) *gorm.DB {
				return tx.Where("id in (select owner_id from role_objectives)")
			}).
		Preload("Factions.Members.RoleObjectives", r.PreloadSelect("OwnerID", "Type", "Completed"))
	processRoots := getRoots(query, c)

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
	processRoots := getRoots(query, c)

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

func HeatmapsGET(c *gin.Context) (int, string, gin.H) {
	return 200, "heatmaps.html", gin.H{}
}

func ChanglingGET(c *gin.Context) (int, string, gin.H) {
	return 200, "changling.html", gin.H{}
}
