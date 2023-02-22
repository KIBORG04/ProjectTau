package domain

import "gorm.io/gorm"

var Models = []any{
	new(Root),
	new(Factions),
	new(Role),
	new(Score),
	new(Rating),
	new(RatingValues),
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
	new(ChangelingInfo),
	new(ChangelingPurchase),
	new(WizardInfo),
	new(WizardPurchase),
	new(Aspects),
	new(RitenameByCount),

	new(Player),
	new(Top),
	new(TopMember),
}

type MyMigrator interface {
	ColumnsMigration(dx *gorm.DB)
}
