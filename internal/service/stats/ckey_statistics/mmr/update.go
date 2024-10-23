package mmr

import (
	"fmt"
	"golang.org/x/exp/slices"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/stats"
)

type playersMMR map[string]int32

func isWin(role *domain.Role, faction *domain.Factions) (int32, error) {
	victoryType, err := stats.GetAntagonistWinType(role.RoleName, faction.FactionName)
	switch victoryType {
	case stats.FactionVictory:
		return faction.Victory, nil
	case stats.RoleVictory:
		return role.Victory, nil
	default:
		return -1, err
	}
}

func ParseMMR() []string {
	query := r.Database.
		Preload("Factions", r.PreloadSelect("RootID", "ID", "FactionName", "Victory")).
		Preload("Factions.Members", r.PreloadSelect("ID", "OwnerID", "MindCkey", "MindName", "RoleName", "Victory")).
		Omit("CompletionHTML")

	var roots []*domain.Root
	query.Find(&roots)

	playersMMR := make(playersMMR)

	for _, root := range roots {
		for _, faction := range root.Factions {
			processedCkeys := make([]string, 0, len(faction.Members))
			for _, role := range faction.Members {
				if role.MindCkey == "" {
					continue
				}
				if slices.Contains(processedCkeys, role.MindCkey) {
					continue
				}

				_, ok := playersMMR[role.MindCkey]
				if !ok {
					playersMMR[role.MindCkey] = 1000
				}

				win, err := isWin(&role, &faction)
				if err != nil {
					continue
				}

				if win == 1 {
					playersMMR[role.MindCkey] = playersMMR[role.MindCkey] + 25
				} else if win == 0 {
					playersMMR[role.MindCkey] = playersMMR[role.MindCkey] - 30
				}
				processedCkeys = append(processedCkeys, role.MindCkey)
			}
		}
	}

	playersSlice := make([]*domain.Player, 0, len(playersMMR))
	for player, mmr := range playersMMR {
		playersSlice = append(playersSlice, &domain.Player{
			Ckey: player,
			MMR:  mmr,
		})
	}
	r.SavePlayers(playersSlice)

	return []string{fmt.Sprintf("For %d players MMR recalculated", len(playersMMR))}
}
