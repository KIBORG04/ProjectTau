package repository

import (
	"fmt"
	"gorm.io/gorm/clause"
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

func EagerFindByRoundId(id string) (*d.Root, error) {
	var root d.Root
	Database.Preload(clause.Associations).
		Preload("Deaths.Damage").
		Preload("Factions").
		Preload("Factions.FactionObjectives").
		Preload("Factions.CultInfo").
		Preload("Factions.Members.RoleObjectives").
		Preload("Factions.Members.UplinkInfo.UplinkPurchases").
		Table("roots").
		Where("round_id = ?", id).First(&root)
	if root.RoundID == 0 {
		return nil, fmt.Errorf("not found %s id", id)
	}
	return &root, nil
}
