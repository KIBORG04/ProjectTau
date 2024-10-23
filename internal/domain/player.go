package domain

import (
	"gorm.io/gorm"
	"time"
)

type Player struct {
	Ckey             string `gorm:"primaryKey;uniqueIndex"`
	MMR              int32
	CrawlerUpdatedAt time.Time
	CrawlerStats     []CrawlerStat
}

func (p *Player) BeforeSave(tx *gorm.DB) (err error) {
	if len(p.CrawlerStats) > 0 {
		p.CrawlerUpdatedAt = time.Now().Truncate(24 * time.Hour)
	}
	return
}

func (p *Player) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(p.CrawlerStats) > 0 {
		p.CrawlerUpdatedAt = time.Now().Truncate(24 * time.Hour)
	}
	return
}

type CrawlerStat struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement"`
	PlayerID         string `gorm:"index"`
	ServerName       string `gorm:"index"`
	CrawlerUpdatedAt time.Time
	Minutes          uint
}

func (p *CrawlerStat) BeforeSave(tx *gorm.DB) (err error) {
	p.CrawlerUpdatedAt = time.Now().Truncate(24 * time.Hour)
	return
}
