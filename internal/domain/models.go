package domain

import "gorm.io/gorm"

var Models = []any{
	new(Root),
	new(Factions),
	new(Role),
	new(Score),
	new(Rating),
	new(RatingValues),
	new(Vote),
	new(VoteValues),
	new(Achievement),
	new(Medal),
	new(CommunicationLogs),
	new(Deaths),
	new(Explosions),
	new(EMPExplosions),
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
	new(CrawlerStat),
	new(Top),
	new(TopMember),

	new(ThetaSponsor),
}

type MyMigrator interface {
	ColumnsMigration(dx *gorm.DB)
}
