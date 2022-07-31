package repository

import "ssstatistics/internal/domain"

func UpdateTop(a []*domain.Top) error {
	err := ClearTopMembers()
	if err != nil {
		return err
	}
	Database.Save(a)
	return nil
}

func ClearTopMembers() error {
	err := Database.Migrator().DropTable(domain.TopMember{})
	if err != nil {
		return err
	}
	err = Database.Migrator().CreateTable(domain.TopMember{})
	if err != nil {
		return err
	}
	return nil
}
