package domain

import "gorm.io/gorm"

var Models = []any{
	new(Root),
	new(Factions),
	new(Role),
	new(Score),
	new(Achievement),
	new(CommunicationLogs),
	new(Deaths),
	new(Explosions),
	new(ManifestEntries),
	new(LeaveStats),
	new(Damage),
	new(RoleObjectives),
	new(FactionObjectives),
	new(CultInfo),
	new(UplinkInfo),
	new(UplinkPurchases),
	new(Aspects),
	new(RitenameByCount),

	new(Player),
	new(MapStats),
	new(MapAttribute),
	new(Top),
	new(TopMember),
}

type MyMigrator interface {
	ColumnsMigration(dx *gorm.DB)
}
