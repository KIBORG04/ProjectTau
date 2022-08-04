package tops

import (
	"golang.org/x/exp/slices"
	"sort"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
	"ssstatistics/internal/utils"
	"strings"
)

type PlayerTopInfo struct {
	Name  string
	Value int
}

type PlayerTopInfoSlice []*PlayerTopInfo

func (info PlayerTopInfoSlice) Len() int           { return len(info) }
func (info PlayerTopInfoSlice) Less(i, j int) bool { return info[i].Value < info[j].Value }
func (info PlayerTopInfoSlice) Swap(i, j int)      { info[i], info[j] = info[j], info[i] }

func (info PlayerTopInfoSlice) hasName(name string) (*PlayerTopInfo, bool) {
	for i := 0; i < len(info); i++ {
		if info[i].Name == name {
			return info[i], true
		}
	}
	return nil, false
}

type TopInfo struct {
	Title           string
	Tag             int
	KeyColumnName   string
	ValueColumnName string
	ValuePostfix    string
	PlayersInfo     PlayerTopInfoSlice
}

func (t *TopInfo) AddPlayerCount(name string) {
	t.ChangePlayerAndValue(name, 1, func(a int, b int) int {
		return a + b
	})
}

func (t *TopInfo) SetPlayerAndMaxValue(name string, value int) {
	t.ChangePlayerAndValue(name, value, func(a int, b int) int {
		return utils.Max(a, b)
	})
}

func (t *TopInfo) SetPlayerAndValue(name string, value int) {
	t.ChangePlayerAndValue(name, value, func(a int, b int) int {
		return b
	})
}

func (t *TopInfo) ChangePlayerAndValue(name string, value int, setRule func(int, int) int) {
	foundInfo, ok := t.PlayersInfo.hasName(name)
	if ok {
		foundInfo.Value = setRule(foundInfo.Value, value)
	} else {
		foundInfo = &PlayerTopInfo{
			Name:  name,
			Value: value,
		}
		t.PlayersInfo = append(t.PlayersInfo, foundInfo)
	}
}

func initStaticTops() map[string]*TopInfo {
	return map[string]*TopInfo{
		"deaths": {
			Title:           "Смертей",
			Tag:             domain.Misc,
			KeyColumnName:   "Имя",
			ValueColumnName: "Количество",
		},
		"zadrots": {
			Title:           "Задротов",
			Tag:             domain.Misc,
			KeyColumnName:   "Имя",
			ValueColumnName: "Раундов",
		},
		"leavers": {
			Title:           "Ливеров",
			Tag:             domain.Misc,
			KeyColumnName:   "Имя",
			ValueColumnName: "Количество",
		},
		"rich": {
			Title:           "Богатейших",
			Tag:             domain.Misc,
			KeyColumnName:   "Имя",
			ValueColumnName: "Денег",
			ValuePostfix:    "$",
		},
		"damaged": {
			Title:           "Избитых",
			Tag:             domain.Misc,
			KeyColumnName:   "Имя",
			ValueColumnName: "Урон",
		},
	}
}

// ParseTopData Deletes the past and adds the current ones
func ParseTopData() []string {
	query := r.Database.
		Preload("Deaths", r.PreloadSelect("RootID", "MindName")).
		Preload("ManifestEntries", r.PreloadSelect("RootID", "AssignedRole", "Name")).
		Preload("LeaveStats", r.PreloadSelect("RootID", "AssignedRole", "Name", "LeaveTime", "LeaveType")).
		Preload("Score", r.PreloadSelect("ID", "RootID", "Richestkey", "Richestcash", "Dmgestkey", "Dmgestdamage")).
		Preload("Factions", r.PreloadSelect("RootID", "ID", "FactionName", "Victory")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "MindCkey", "MindName", "RoleName", "Victory")).
		Omit("CompletionHTML")

	var roots []*domain.Root
	query.Find(&roots)

	staticTopTypes := initStaticTops()
	gamemodeTops := make(map[string]*TopInfo)

	hasId := func(topMap map[string]*TopInfo, id string) bool {
		_, ok := topMap[id]
		return ok
	}
	getTopById := func(topMap map[string]*TopInfo, id string) *TopInfo {
		return topMap[id]
	}

	ckeyAntagsPlayed := make(map[string]map[string]map[string]uint)
	for _, root := range roots {
		for _, death := range root.Deaths {
			// wtf
			if death.MindName == "Unknown" || death.MindName == "unknown" {
				continue
			}
			staticTopTypes["deaths"].AddPlayerCount(death.MindName)
		}
		for _, entry := range root.ManifestEntries {
			if stats.IsStationPlayer(entry.AssignedRole, entry.Name) {
				staticTopTypes["zadrots"].AddPlayerCount(entry.Name)
			}
		}
		for _, stat := range root.LeaveStats {
			if stats.IsStationPlayer(stat.AssignedRole, stat.Name) && stats.IsRoundStartLeaver(stat) {
				staticTopTypes["leavers"].AddPlayerCount(stat.Name)
			}
		}

		for _, faction := range root.Factions {
			if slices.Contains(stats.TeamlRoles, faction.FactionName) && !hasId(gamemodeTops, stats.Ckey(faction.FactionName)) {
				title := faction.FactionName
				if value, ok := stats.ShortModeName[faction.FactionName]; ok {
					title = value
				}
				gamemodeTops[stats.Ckey(faction.FactionName)] = &TopInfo{
					Title:           title,
					Tag:             domain.Gamemode,
					KeyColumnName:   "Ckey",
					ValueColumnName: "Winrate",
					ValuePostfix:    "%",
				}
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

				if slices.Contains(stats.SoloRoles, role.RoleName) && !hasId(gamemodeTops, stats.Ckey(role.RoleName)) {
					title := role.RoleName
					if value, ok := stats.ShortModeName[role.RoleName]; ok {
						title = value
					}
					gamemodeTops[stats.Ckey(role.RoleName)] = &TopInfo{
						Title:           title,
						Tag:             domain.Gamemode,
						KeyColumnName:   "Ckey",
						ValueColumnName: "Winrate",
						ValuePostfix:    "%",
					}
				}

				if _, ok := ckeyAntagsPlayed[role.MindCkey]; !ok {
					antagMap := make(map[string]map[string]uint)
					ckeyAntagsPlayed[role.MindCkey] = antagMap
				}

				if slices.Contains(stats.TeamlRoles, faction.FactionName) {
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

				if slices.Contains(stats.SoloRoles, role.RoleName) {
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

		if root.Score.Richestkey != "" {
			staticTopTypes["rich"].SetPlayerAndMaxValue(root.Score.Richestkey, int(root.Score.Richestcash))
		}
		if root.Score.Dmgestkey != "" {
			staticTopTypes["damaged"].SetPlayerAndMaxValue(root.Score.Dmgestkey, int(root.Score.Dmgestdamage))
		}
	}

	for player, antagInfo := range ckeyAntagsPlayed {
		for antag, antagOptions := range antagInfo {
			if antagOptions["Count"] > 10 {
				antagTop := getTopById(gamemodeTops, stats.Ckey(antag))
				antagTop.SetPlayerAndValue(player, int(float32(antagOptions["Victory"]*100)/float32(antagOptions["Count"])))
			}
		}
	}

	generalTopSlice := make(map[string]*TopInfo, len(gamemodeTops)+len(staticTopTypes))
	for key, top := range staticTopTypes {
		generalTopSlice[key] = top
	}
	for key, top := range gamemodeTops {
		generalTopSlice[key] = top
	}
	// remove useless positions
	for key, top := range generalTopSlice {
		if len(top.PlayersInfo) == 0 {
			delete(generalTopSlice, key)
		}
		sort.Sort(sort.Reverse(top.PlayersInfo))
		if len(top.PlayersInfo) > 10 {
			top.PlayersInfo = slices.Delete(top.PlayersInfo, 10, len(top.PlayersInfo))
		}
	}

	databaseTops := make([]*domain.Top, 0, len(generalTopSlice))
	for id, info := range generalTopSlice {
		members := make([]domain.TopMember, 0, len(info.PlayersInfo))
		for _, playerInfo := range info.PlayersInfo {
			members = append(members, domain.TopMember{
				Key:   playerInfo.Name,
				Value: playerInfo.Value,
			})
		}
		databaseTops = append(databaseTops, &domain.Top{
			ID:              id,
			Title:           info.Title,
			Tag:             info.Tag,
			KeyColumnName:   info.KeyColumnName,
			ValueColumnName: info.ValueColumnName,
			ValuePostfix:    info.ValuePostfix,
			TopMembers:      members,
		})
	}

	err := r.UpdateTop(databaseTops)
	if err != nil {
		return []string{err.Error()}
	}
	return []string{}
}
