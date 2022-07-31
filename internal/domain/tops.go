package domain

// tags
const (
	Misc int = iota
	Gamemode
)

type Top struct {
	ID              string `gorm:"primarykey"`
	Title           string `gorm:"size:256"`
	Tag             int
	KeyColumnName   string `gorm:"size:256"`
	ValueColumnName string `gorm:"size:256"`
	ValuePostfix    string `gorm:"size:32"`

	TopMembers []TopMember
}

type TopMember struct {
	ID    uint `gorm:"primarykey"`
	TopID string
	Key   string `gorm:"size:256"`
	Value int
}
