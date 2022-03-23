package domain

import "gorm.io/gorm"

var Models = []any{
	&Root{},
	&Factions{},
	&Role{},
	&Score{},
	&Achievement{},
	&CommunicationLogs{},
	&Deaths{},
	&Explosions{},
	&ManifestEntries{},
	&LeaveStats{},
	&Damage{},
	&Objectives{},
	&CultInfo{},
	&UplinkInfo{},
	&UplinkPurchases{},
	&Aspects{},
	&RitenameByCount{},
	&Link{},
}

type MyMigrator interface {
	ColumnsMigration(dx *gorm.DB)
}
