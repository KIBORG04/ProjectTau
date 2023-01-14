package general

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ssstatistics/internal/config"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/utils"
)

func BasicGET(c *gin.Context, f func(*gin.Context) (int, string, gin.H)) {
	baseContext := gin.H{
		"baseUrl":          config.Config.BaseUrl,
		"serverCheckboxes": stats.GetCheckboxStates(c),
	}

	code, file, context := f(c)

	for k, v := range baseContext {
		context[k] = v
	}

	c.HTML(code, file, context)
}

func BasicPOST(c *gin.Context, f func(*gin.Context) (int, string, gin.H)) {
	stats.SetCheckboxStates(c)
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
		Where("assigned_role in (?)", stats.StationPositions).
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
		case stats.ServerAlphaAddress:
			totalAlphaRoots += root.Count
		case stats.ServerBetaAddress:
			totalBetaaRoots += root.Count
		case stats.ServerGammaAddress:
			totalGammaRoots += root.Count
		}
	}

	return 200, "index.html", gin.H{
		"totalRounds": totalRoots,
		"version":     lastRoot.Version,
		"lastRound":   lastRoot.RoundID,
		"lastDate":    utils.TrimPGDate(lastRoot.Date),
		"firstDate":   stats.CurrentStatisticsDate,

		"alphaRounds": totalAlphaRoots,
		"betaRounds":  totalBetaaRoots,
		"gammaRounds": totalGammaRoots,

		"modesCount":      modesCount,
		"modesSum":        modesSum,
		"crewDeathsCount": crewDeathsCount,
		"crewDeathsSum":   crewDeathsSum,
		"roleDeathsCount": roleDeathsCount,
		"roleDeathsSum":   roleDeathsSum,
	}
}

func GamemodesGET(c *gin.Context) (int, string, gin.H) {
	var (
		factionsStatistics []*struct {
			FactionName           string
			AvgLeavers            float32
			Count                 int
			Wins                  int
			Winrate               float32
			WinrateUint           uint
			MembersCount          int
			TotalObjectives       int
			CompletedObjectives   int
			WinrateObjectives     float32
			WinrateObjectivesUint uint
		}

		rolesStatistics []*struct {
			RoleName              string
			Count                 int
			Wins                  int
			Winrate               float32
			WinrateUint           uint
			TotalObjectives       int
			CompletedObjectives   int
			WinrateObjectives     float32
			WinrateObjectivesUint uint
		}
	)

	r.Database.Raw(`
SELECT factions.faction_name,
       (select (select sum(a.leaves)
                from (SELECT distinct on (leave_stats.id) CASE
                              WHEN leave_time = ''
                                  THEN 0
                              WHEN leave_type = 'Cryopod' AND split_part(leave_time, ':', 2)::int < 15
                                  THEN 0
                              WHEN split_part(leave_time, ':', 2)::int > 5 AND split_part(leave_time, ':', 2)::int < 30
                                  THEN 1
                              WHEN leave_type = 'Cryopod' AND split_part(leave_time, ':', 2)::int < 45
                                  THEN 1
                              END AS leaves
                      FROM leave_stats
                      JOIN roots aaa on aaa.round_id = leave_stats.root_id
                      JOIN factions bbb on bbb.root_id = aaa.round_id
                      WHERE bbb.faction_name = factions.faction_name
                        AND assigned_role IN ?
                        AND leave_stats.name NOT LIKE 'maintenance drone%') as a))::real / COUNT(factions.id)                                                          AS avg_leavers,
       COUNT(factions.id)                                                                                                                                              AS count,
       SUM(victory)                                                                                                                                                    AS wins,
       SUM(victory)::real * 100 / COUNT(factions.id)::real                                                                                                                   AS winrate,
       SUM((SELECT count(1) FROM roles where roles.owner_id = factions.id))                                                     AS members_count,
       SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id))                                      AS total_objectives,
       SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id and fo1.completed = 'SUCCESS'))        AS completed_objectives,
       SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id and fo1.completed = 'SUCCESS'))::real * 100 /
            GREATEST(SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id)) ::real, 1)             AS winrate_objectives

		FROM factions
		group by factions.faction_name;
	`, stats.StationPositions).Scan(&factionsStatistics)

	r.Database.Raw(`
	select roles.role_name,
       COUNT(1)                                                                                                  AS count,
       SUM(roles.victory)                                                                                        AS wins,
       SUM(roles.victory)::real * 100 / COUNT(1)::real                                                           AS winrate,
       SUM((SELECT count(1)
            FROM role_objectives
            WHERE roles.id = role_objectives.owner_id))                                                          AS total_objectives,
       SUM((SELECT count(1)
            FROM role_objectives
            WHERE roles.id = role_objectives.owner_id
              AND completed = 'SUCCESS'))                                                                        AS completed_objectives,
       SUM((SELECT count(1)
            FROM role_objectives
            WHERE roles.id = role_objectives.owner_id
              AND completed = 'SUCCESS'))::real * 100 /
       GREATEST(SUM((SELECT count(1) FROM role_objectives WHERE roles.id = role_objectives.owner_id)) ::real,
                1)                                                                                               AS winrate_objectives

	from roles
	group by roles.role_name;
	`).Scan(&rolesStatistics)

	var factionsSum int
	var rolesSum int

	for _, statistic := range factionsStatistics {
		factionsSum += statistic.Count
		statistic.WinrateUint = uint(statistic.Winrate)
		statistic.WinrateObjectivesUint = uint(statistic.WinrateObjectives)
	}
	for _, statistic := range rolesStatistics {
		rolesSum += statistic.Count
		statistic.WinrateUint = uint(statistic.Winrate)
		statistic.WinrateObjectivesUint = uint(statistic.WinrateObjectives)
	}

	return 200, "gamemodes.html", gin.H{
		"factionsSum":   factionsSum,
		"factionsCount": factionsStatistics,
		"rolesSum":      rolesSum,
		"rolesCount":    rolesStatistics,
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
	ObjectiveInfos stats.InfoSlice
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
	processRoots := stats.GetRoots(query, c)

	objectiveHolders := make(stats.InfoSlice, 0)

	addObjectiveInfo := func(owner *OwnerByObjectivesInfo, objective domain.Objectives) {
		owner.Count++
		var isWin uint
		if objective.Completed == stats.ObjectiveWIN {
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
						Id:    stats.Ckey(faction.FactionName),
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
							Id:    stats.Ckey(role.RoleName),
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
	processRoots := stats.GetRoots(query, c)

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
