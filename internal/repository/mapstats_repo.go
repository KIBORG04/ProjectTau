package repository

import "ssstatistics/internal/domain"

func ClearMapStats() error {
	err := Database.Migrator().DropTable(domain.MapStats{})
	if err != nil {
		return err
	}
	err = Database.Migrator().CreateTable(domain.MapStats{})
	if err != nil {
		return err
	}
	err = Database.Migrator().DropTable(domain.MapAttribute{})
	if err != nil {
		return err
	}
	err = Database.Migrator().CreateTable(domain.MapAttribute{})
	if err != nil {
		return err
	}
	return nil
}

func SaveMapStats(maps []*domain.MapStats) {
	err := ClearMapStats()
	if err != nil {
		return
	}
	Database.Save(maps)
}
