package maps

import (
	"fmt"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

type mapAttrubite struct {
	Name  string
	Sum   float32
	Total uint
}

type mapStruct struct {
	MapName  string
	ServerID string

	MapAttributes []*mapAttrubite
}

func FixMaxShit() []string {
	r.Database.Model(new(domain.Root)).Where("map = ?", "Prometheus").Update("map", "Prometheus Station")
	return []string{"Map Shit Fixed"}
}

var maps []*mapStruct

func findMapByNameAndServerID(name, id string) (*mapStruct, bool) {
	for _, mapStat := range maps {
		if mapStat.MapName == name && mapStat.ServerID == id {
			return mapStat, true
		}
	}
	return nil, false
}

func findMapAttributeByName(mapStat *mapStruct, name string) (*mapAttrubite, bool) {
	for _, attribute := range mapStat.MapAttributes {
		if attribute.Name == name {
			return attribute, true
		}
	}
	return nil, false
}

func addAvgMapAttribute(mapStat *mapStruct, name string, value float32) {
	attribute, ok := findMapAttributeByName(mapStat, name)
	if !ok {
		attribute = &mapAttrubite{
			Name:  name,
			Sum:   value,
			Total: 1,
		}
		mapStat.MapAttributes = append(mapStat.MapAttributes, attribute)
	}

	attribute.Sum += value
	attribute.Total += 1
}

func addMapAttribute(mapStat *mapStruct, name string, value float32) {
	attribute, ok := findMapAttributeByName(mapStat, name)
	if !ok {
		attribute = &mapAttrubite{
			Name:  name,
			Sum:   value,
			Total: 1,
		}
		mapStat.MapAttributes = append(mapStat.MapAttributes, attribute)
	}

	attribute.Sum += value
}

func ParseMapsInfo() []string {
	query := r.Database.
		Preload("Score").
		Omit("CompletionHTML", "TestMerges", "BaseCommitSha", "ModeResult")

	var roots []*domain.Root
	query.Find(&roots)

	for _, root := range roots {
		mapStat, ok := findMapByNameAndServerID(root.Map, root.ServerAddress)
		if !ok {
			mapStat = &mapStruct{
				MapName:  root.Map,
				ServerID: root.ServerAddress,
			}
			maps = append(maps, mapStat)
		}

		roundTime, err := stats.ParseRoundTime(root.Duration)
		if err == nil {
			addAvgMapAttribute(mapStat, "AvgDuration", float32(roundTime.ToSeconds()))
		} else {
			addAvgMapAttribute(mapStat, "AvgDuration", float32(3600))
		}

		addMapAttribute(mapStat, "Picked", 1)

		addAvgMapAttribute(mapStat, "AvgCrewsocre", float32(root.Score.Crewscore))
		addAvgMapAttribute(mapStat, "AvgStuffshipped", float32(root.Score.Stuffshipped))
		addAvgMapAttribute(mapStat, "AvgStuffharvested", float32(root.Score.Stuffharvested))
		addAvgMapAttribute(mapStat, "AvgOremined", float32(root.Score.Oremined))
		addAvgMapAttribute(mapStat, "AvgResearchdone", float32(root.Score.Researchdone))
		addAvgMapAttribute(mapStat, "AvgPowerloss", float32(root.Score.Powerloss))
		addAvgMapAttribute(mapStat, "AvgMess", float32(root.Score.Mess))
		addAvgMapAttribute(mapStat, "AvgMeals", float32(root.Score.Meals))
		addAvgMapAttribute(mapStat, "AvgNuked", float32(root.Score.Nuked))
		addAvgMapAttribute(mapStat, "AvgRecAntags", float32(root.Score.RecAntags))
		addAvgMapAttribute(mapStat, "AvgCrewEscaped", float32(root.Score.CrewEscaped))
		addAvgMapAttribute(mapStat, "AvgCrewDead", float32(root.Score.CrewDead))
		addAvgMapAttribute(mapStat, "AvgCrewTotal", float32(root.Score.CrewTotal))
		addAvgMapAttribute(mapStat, "AvgCrewSurvived", float32(root.Score.CrewSurvived))
		addAvgMapAttribute(mapStat, "AvgFoodeaten", float32(root.Score.Foodeaten))
		addAvgMapAttribute(mapStat, "AvgClownabuse", float32(root.Score.Clownabuse))

	}

	mapStats := make([]*domain.MapStats, 0, len(maps))
	for _, m := range maps {
		mapStat := &domain.MapStats{
			MapName:  m.MapName,
			ServerID: m.ServerID,
		}

		for _, attribute := range m.MapAttributes {
			mapStat.MapAttributes = append(mapStat.MapAttributes, &domain.MapAttribute{
				Name:  attribute.Name,
				Value: fmt.Sprint(attribute.Sum / float32(attribute.Total)),
			})
		}

		mapStats = append(mapStats, mapStat)
	}

	r.SaveMapStats(mapStats)

	return []string{fmt.Sprintf("For %d map statistics parsed", len(mapStats))}
}
