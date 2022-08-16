package domain

type Player struct {
	Ckey string `gorm:"primaryKey;uniqueIndex"`
	MMR  int32
}
