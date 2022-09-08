package domain

type MapAttribute struct {
	ID         uint `gorm:"primarykey"`
	MapStatsID string

	Name  string
	Value string
}

type MapStats struct {
	ID       uint `gorm:"primarykey"`
	MapName  string
	ServerID string

	MapAttributes []*MapAttribute
}
