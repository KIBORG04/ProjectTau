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
	&RoleObjectives{},
	&FactionObjectives{},
	&CultInfo{},
	&UplinkInfo{},
	&UplinkPurchases{},
	&Aspects{},
	&RitenameByCount{},

	&Player{},
}

type MyMigrator interface {
	ColumnsMigration(dx *gorm.DB)
}
