package repository

import "ssstatistics/internal/domain"

func SaveMMR(players []*domain.Player) {
	Database.Save(players)
}
