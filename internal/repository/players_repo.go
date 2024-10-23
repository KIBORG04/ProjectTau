package repository

import (
	"fmt"
	"ssstatistics/internal/domain"
)

func SavePlayers(players []*domain.Player) {
	Database.Save(players)
}

func SavePlayer(player *domain.Player) {
	Database.Save(player)
}

func GetPlayer(ckey string) *domain.Player {
	var player domain.Player
	result := Database.Preload("CrawlerStats").Where("ckey = ?", ckey).First(&player)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil
	}
	return &player
}
