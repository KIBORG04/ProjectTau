package repository

import (
	"fmt"
	d "ssstatistics/internal/domain"
)

func FindByRoundId(id string) (*d.Root, error) {
	var root d.Root
	Database.Table("roots").Where("round_id = ?", id).First(&root)
	if root.RoundID == 0 {
		return nil, fmt.Errorf("not found %s id", id)
	}
	return &root, nil
}
